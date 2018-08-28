package pefile

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type PE struct {
	Type                 OptionalHeaderMagic
	DosHeader            DosHeader
	NtHeaders            NtHeaders
	FileHeader           FileHeader
	OptionalHeader       *OptionalHeader
	OptionalHeader64     *OptionalHeader64
	RichHeader           *RichHeader
	DataDirectoryEntries []DataDirectoryEntry
	DataDirectory        map[DirectoryEntryType]interface{}
	Sections             []SectionHeader

	Warnings []string
}

func (p *PE) warn(format string, args ...interface{}) {
	p.Warnings = append(p.Warnings, fmt.Sprintf(format, args...))
}

func Parse(buf []byte) (*PE, error) {
	p := &PE{}
	r := bytes.NewReader(buf)
	if err := binary.Read(r, binary.LittleEndian, &p.DosHeader); err != nil {
		return nil, err
	}
	switch p.DosHeader.E_magic {
	case IMAGE_DOSZM_SIGNATURE:
		return nil, errors.New("Probably a ZM Executable (not a PE file)")
	case IMAGE_DOS_SIGNATURE:
	default:
		return nil, errors.New("DOS Header magic not found")
	}
	if int(p.DosHeader.E_lfanew) > len(buf) {
		return nil, errors.New("Invalid e_lfanew value, probably not a PE file")
	}

	// pefile.py:1843
	ntHeadersOffset := p.DosHeader.E_lfanew
	if _, err := r.Seek(int64(ntHeadersOffset), 0); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &p.NtHeaders); err != nil {
		return nil, err
	}
	switch 0xffff & p.NtHeaders.Signature {
	case IMAGE_NE_SIGNATURE:
		return nil, errors.New("Invalid NT Headers signature. Probably a NE file")
	case IMAGE_LE_SIGNATURE:
		return nil, errors.New("Invalid NT Headers signature. Probably a LE file")
	case IMAGE_LX_SIGNATURE:
		return nil, errors.New("Invalid NT Headers signature. Probably a LX file")
	case IMAGE_TE_SIGNATURE:
		return nil, errors.New("Invalid NT Headers signature. Probably a TE file")
	}
	if p.NtHeaders.Signature != IMAGE_NT_SIGNATURE {
		return nil, errors.New("Invalid NT Headers signature")
	}

	// pefile.py:1866
	if _, err := r.Seek(int64(ntHeadersOffset)+4, 0); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &p.FileHeader); err != nil {
		return nil, err
	}
	// TODO ASCII representation?
	// image_flags = retrieve_flags(IMAGE_CHARACTERISTICS, 'IMAGE_FILE_')

	// sizeof FileHeader == 20
	optionalHeaderOffset := ntHeadersOffset + 4 + 20
	sectionsOffset := optionalHeaderOffset + uint32(p.FileHeader.SizeOfOptionalHeader)

	// pefile.py:1885
	if _, err := r.Seek(int64(optionalHeaderOffset), 0); err != nil {
		return nil, err
	}
	var ohbuf [96]byte
	if n, _ := r.Read(ohbuf[:]); n >= 69 {
		ohr := bytes.NewReader(ohbuf[:])
		var optionalHeader OptionalHeader
		if err := binary.Read(ohr, binary.LittleEndian, &optionalHeader); err != nil {
			return nil, err
		}
		p.OptionalHeader = &optionalHeader
	}

	// pefile.py:1923
	if p.OptionalHeader != nil {
		switch p.OptionalHeader.Magic {
		case OPTIONAL_HEADER_MAGIC_PE:
			p.Type = OPTIONAL_HEADER_MAGIC_PE
		case OPTIONAL_HEADER_MAGIC_PE_PLUS:
			p.Type = OPTIONAL_HEADER_MAGIC_PE_PLUS
			if _, err := r.Seek(int64(optionalHeaderOffset), 0); err != nil {
				return nil, err
			}
			var ohbuf [112]byte
			if n, _ := r.Read(ohbuf[:]); n >= 73 {
				ohr := bytes.NewReader(ohbuf[:])
				var optionalHeader OptionalHeader64
				if err := binary.Read(ohr, binary.LittleEndian, &optionalHeader); err != nil {
					return nil, err
				}
				p.OptionalHeader64 = &optionalHeader
			}
		}
	}

	if p.Type == 0 || p.OptionalHeader == nil && p.OptionalHeader64 == nil {
		return nil, errors.New("No Optional Header found, invalid PE32 or PE32+ file")
	}

	// TODO, see above
	// dll_characteristics_flags = retrieve_flags(DLL_CHARACTERISTICS, 'IMAGE_DLLCHARACTERISTICS_')

	var (
		addressOfEntryPoint uint32
		sizeOfHeaders       uint32
		numberOfRvaAndSizes uint32
		fileAlignment       uint32
	)
	switch true {
	case p.OptionalHeader != nil:
		addressOfEntryPoint = p.OptionalHeader.AddressOfEntryPoint
		sizeOfHeaders = p.OptionalHeader.SizeOfHeaders
		numberOfRvaAndSizes = p.OptionalHeader.NumberOfRvaAndSizes
		fileAlignment = p.OptionalHeader.FileAlignment
	case p.OptionalHeader64 != nil:
		addressOfEntryPoint = p.OptionalHeader64.AddressOfEntryPoint
		sizeOfHeaders = p.OptionalHeader64.SizeOfHeaders
		numberOfRvaAndSizes = p.OptionalHeader64.NumberOfRvaAndSizes
		fileAlignment = p.OptionalHeader64.FileAlignment
	}
	if addressOfEntryPoint < sizeOfHeaders {
		p.warn("SizeOfHeaders is smaller than AddressOfEntryPoint: this file cannot run under Windows 8")
	}
	if numberOfRvaAndSizes > 0x10 {
		p.warn("Suspicious NumberOfRvaAndSizes in the Optional Header."+
			"Normal values are never larger than 0x10, the value is: 0x%x", numberOfRvaAndSizes)
	}

	offset := optionalHeaderOffset + uint32(p.FileHeader.SizeOfOptionalHeader)

	// pefile.py:2000
	for i := 0; i < int(0x7fffffff&numberOfRvaAndSizes); i++ {
		if offset == uint32(len(buf)) {
			break
		}
		if _, err := r.Seek(int64(optionalHeaderOffset), 0); err != nil {
			return nil, err
		}
		// TODO: HOWTO: dir_entry.name = DIRECTORY_ENTRY[i]
		var dirEntry DataDirectoryEntry
		if err := binary.Read(r, binary.LittleEndian, &dirEntry); err != nil {
			return nil, err
		}
		p.DataDirectoryEntries = append(p.DataDirectoryEntries, dirEntry)
		offset += 8
		// TODO: Workarounds for broken files
	}

	// pefile.py:2045
	// offset = self.parse_sections(sections_offset)
	offset = sectionsOffset
	for i := 0; i < int(p.FileHeader.NumberOfSections); i++ {
		var sectionHeader SectionHeader
		var nullHeader SectionHeader
		if _, err := r.Seek(int64(offset), 0); err != nil {
			return nil, err
		}
		if err := binary.Read(r, binary.LittleEndian, &sectionHeader); err != nil {
			return nil, err
		}
		if sectionHeader == nullHeader {
			p.warn("Invalid section %d. Contents are null-bytes.", i)
		}
		// TODO: is this corkami\'s virtsectblXP?
		if int(sectionHeader.SizeOfRawData) > len(buf) {
			p.warn("Error parsing section %d. SizeOfRawData is larger than file.", i)
		}
		// TODO: if self.adjust_FileAlignment( section.PointerToRawData,
		if sectionHeader.Misc > 0x10000000 {
			p.warn("Suspicious value found parsing section %d. VirtualSize is extremely large > 256MiB.", i)
		}
		// TODO: if self.adjust_SectionAlignment( section.VirtualAddress,
		if fileAlignment != 0 && (sectionHeader.PointerToRawData%fileAlignment) != 0 {
			p.warn("Error parsing section %d. "+
				"PointerToRawData should normally be "+
				"a multiple of FileAlignment, this might imply the file "+
				"is trying to confuse tools which parse this incorrectly.", i)
		}
		// TODO: PAGE, driver... (pefile.py:2383)
		// TODO: detect overlapping sections (pefile.py:2399)

		p.Sections = append(p.Sections, sectionHeader)

		offset += 40 // sizeof SectionHeader
	}
	// TODO: Check whether the entry point lives within a section / within the file (pefile.py:2070)

	// pefile.py:2093
	// self.parse_data_directories()
	// pefile.py:2416
	for dirEntryType, _ := range p.DataDirectoryEntries {
		switch DirectoryEntryType(dirEntryType) {
		case IMAGE_DIRECTORY_ENTRY_EXPORT:
		case IMAGE_DIRECTORY_ENTRY_IMPORT:
		case IMAGE_DIRECTORY_ENTRY_RESOURCE:
		case IMAGE_DIRECTORY_ENTRY_EXCEPTION:
		case IMAGE_DIRECTORY_ENTRY_SECURITY:
		case IMAGE_DIRECTORY_ENTRY_BASERELOC:
		case IMAGE_DIRECTORY_ENTRY_DEBUG:
		case IMAGE_DIRECTORY_ENTRY_COPYRIGHT:
		case IMAGE_DIRECTORY_ENTRY_GLOBALPTR:
		case IMAGE_DIRECTORY_ENTRY_TLS:
		case IMAGE_DIRECTORY_ENTRY_LOAD_CONFIG:
		case IMAGE_DIRECTORY_ENTRY_BOUND_IMPORT:
		case IMAGE_DIRECTORY_ENTRY_IAT:
		case IMAGE_DIRECTORY_ENTRY_DELAY_IMPORT:
		case IMAGE_DIRECTORY_ENTRY_COM_DESCRIPTOR:
		case IMAGE_DIRECTORY_ENTRY_RESERVED:
		}
	}

	// pe.file.py:2097
	// rich_header = self.parse_rich_header()

	return p, nil
}

func Load(filename string) (*PE, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	f.Close()
	return Parse(buf)
}
