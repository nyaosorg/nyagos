package pefile

type DosHeader struct {
	E_magic    uint16
	E_cblp     uint16
	E_cp       uint16
	E_crlc     uint16
	E_cparhdr  uint16
	E_minalloc uint16
	E_maxalloc uint16
	E_ss       uint16
	E_sp       uint16
	E_csum     uint16
	E_ip       uint16
	E_cs       uint16
	E_lfarlc   uint16
	E_ovno     uint16
	E_res      [8]byte
	E_oemid    uint16
	E_oeminfo  uint16
	E_res2     [20]byte
	E_lfanew   uint32
}

type FileHeader struct {
	Machine              MachineType
	NumberOfSections     uint16
	TimeDateStamp        TimeDateStamp
	PointerToSymbolTable uint32
	NumberOfSymbols      uint32
	SizeOfOptionalHeader uint16
	Characteristics      ImageCharacteristics
}

type NtHeaders struct {
	Signature uint32
}

type OptionalHeader struct {
	Magic                       OptionalHeaderMagic
	MajorLinkerVersion          uint8
	MinorLinkerVersion          uint8
	SizeOfCode                  uint32
	SizeOfInitializedData       uint32
	SizeOfUninitializedData     uint32
	AddressOfEntryPoint         uint32
	BaseOfCode                  uint32
	BaseOfData                  uint32
	ImageBase                   uint32
	SectionAlignment            uint32
	FileAlignment               uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion           uint16
	MinorImageVersion           uint16
	MajorSubsystemVersion       uint16
	MinorSubsystemVersion       uint16
	Reserved1                   uint32
	SizeOfImage                 uint32
	SizeOfHeaders               uint32
	CheckSum                    uint32
	Subsystem                   SubsystemType
	DllCharacteristics          DllCharacteristics
	SizeOfStackReserve          uint32
	SizeOfStackCommit           uint32
	SizeOfHeapReserve           uint32
	SizeOfHeapCommit            uint32
	LoaderFlags                 uint32
	NumberOfRvaAndSizes         uint32
}

type OptionalHeader64 struct {
	Magic                       OptionalHeaderMagic
	MajorLinkerVersion          uint8
	MinorLinkerVersion          uint8
	SizeOfCode                  uint32
	SizeOfInitializedData       uint32
	SizeOfUninitializedData     uint32
	AddressOfEntryPoint         uint32
	BaseOfCode                  uint32
	ImageBase                   uint64
	SectionAlignment            uint32
	FileAlignment               uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion           uint16
	MinorImageVersion           uint16
	MajorSubsystemVersion       uint16
	MinorSubsystemVersion       uint16
	Reserved1                   uint32
	SizeOfImage                 uint32
	SizeOfHeaders               uint32
	CheckSum                    uint32
	Subsystem                   uint16
	DllCharacteristics          DllCharacteristics
	SizeOfStackReserve          uint64
	SizeOfStackCommit           uint64
	SizeOfHeapReserve           uint64
	SizeOfHeapCommit            uint64
	LoaderFlags                 uint32
	NumberOfRvaAndSizes         uint32
}

type DirectoryEntryType uint16

const (
	IMAGE_DIRECTORY_ENTRY_EXPORT         DirectoryEntryType = 0
	IMAGE_DIRECTORY_ENTRY_IMPORT         DirectoryEntryType = 1
	IMAGE_DIRECTORY_ENTRY_RESOURCE       DirectoryEntryType = 2
	IMAGE_DIRECTORY_ENTRY_EXCEPTION      DirectoryEntryType = 3
	IMAGE_DIRECTORY_ENTRY_SECURITY       DirectoryEntryType = 4
	IMAGE_DIRECTORY_ENTRY_BASERELOC      DirectoryEntryType = 5
	IMAGE_DIRECTORY_ENTRY_DEBUG          DirectoryEntryType = 6
	IMAGE_DIRECTORY_ENTRY_COPYRIGHT      DirectoryEntryType = 7 // Architecture on non-x86 platforms
	IMAGE_DIRECTORY_ENTRY_GLOBALPTR      DirectoryEntryType = 8
	IMAGE_DIRECTORY_ENTRY_TLS            DirectoryEntryType = 9
	IMAGE_DIRECTORY_ENTRY_LOAD_CONFIG    DirectoryEntryType = 10
	IMAGE_DIRECTORY_ENTRY_BOUND_IMPORT   DirectoryEntryType = 11
	IMAGE_DIRECTORY_ENTRY_IAT            DirectoryEntryType = 12
	IMAGE_DIRECTORY_ENTRY_DELAY_IMPORT   DirectoryEntryType = 13
	IMAGE_DIRECTORY_ENTRY_COM_DESCRIPTOR DirectoryEntryType = 14
	IMAGE_DIRECTORY_ENTRY_RESERVED       DirectoryEntryType = 15
)

type DataDirectoryEntry struct {
	VirtualAddress uint32
	Size           uint32
}

type SectionHeader struct {
	Name [8]byte
	// union
	// {
	//     ULONG PhysicalAddress;
	//     ULONG VirtualSize;
	// } Misc;
	Misc                 uint32
	VirtualAddress       uint32
	SizeOfRawData        uint32
	PointerToRawData     uint32
	PointerToRelocations uint32
	PointerToLinenumbers uint32
	NumberOfRelocations  uint16
	NumberOfLinenumbers  uint16
	Characteristics      SectionCharacteristics
}

type DelayImportDescriptor struct {
	grAttrs     uint32
	szName      uint32
	phmod       uint32
	pIAT        uint32
	pINT        uint32
	pBoundIAT   uint32
	pUnloadIAT  uint32
	dwTimeStamp TimeDateStamp
}

type ImportDescriptor struct {
	Characteristics ImportDescriptorCharacteristics // union with OriginalFirstChunk
	TimeDateStamp   TimeDateStamp
	ForwarderChain  uint32
	Name            uint32
	FirstThunk      uint32
}

type ExportDirectory struct {
	Characteristics       ExportCharacteristics
	TimeDateStamp         TimeDateStamp
	MajorVersion          uint16
	MinorVersion          uint16
	Name                  uint32
	Base                  uint32
	NumberOfFunctions     uint32
	NumberOfNames         uint32
	AddressOfFunctions    uint32
	AddressOfNames        uint32
	AddressOfNameOrdinals uint32
}

type ResourceDirecotry struct {
	Characteristics      ResourceDirecotryCharacteristics
	TimeDateStamp        TimeDateStamp
	MajorVersion         uint16
	MinorVersion         uint16
	NumberOfNamedEntries uint16
	NumberOfIdEntries    uint16
}

type ResourceDirectoryEntry struct {
	Name         uint32
	OffsetToData uint32
}

type ResourceDataEntry struct {
	OffsetToData uint32
	Size         uint32
	CodePage     uint32
	Reserved     uint32
}

type VersionInfo struct {
	Length      uint16
	ValueLength uint16
	Type        uint16
}

type FixedFileInfo struct {
	Signature        uint32
	StrucVersion     uint32
	FileVersionMS    uint32
	FileVersionLS    uint32
	ProductVersionMS uint32
	ProductVersionLS uint32
	FileFlagsMask    uint32
	FileFlags        uint32
	FileOS           uint32
	FileType         uint32
	FileSubtype      uint32
	FileDateMS       uint32
	FileDateLS       uint32
}

type ThunkData uint32

// typedef struct _IMAGE_THUNK_DATA32 {
//   union {
// DWORD ForwarderString;
// DWORD Function;
// DWORD Ordinal;
// DWORD AddressOfData;
//   } u1;
// } IMAGE_THUNK_DATA32;

type ThunkData64 uint64

// typedef struct _IMAGE_THUNK_DATA64 {
//   union {
// ULONGLONG ForwarderString;
// ULONGLONG Function;
// ULONGLONG Ordinal;
// ULONGLONG AddressOfData;
//   } u1;
// } IMAGE_THUNK_DATA64;

type DebugDirectory struct {
	Characteristics  DebugDirectoryCharacteristics
	TimeDateStamp    TimeDateStamp
	MajorVersion     uint16
	MinorVersion     uint16
	Type             uint32
	SizeOfData       uint32
	AddressOfRawData uint32
	PointerToRawData uint32
}

type RelocationType uint16

const (
	IMAGE_REL_BASED_ABSOLUTE       RelocationType = 0
	IMAGE_REL_BASED_HIGH           RelocationType = 1
	IMAGE_REL_BASED_LOW            RelocationType = 2
	IMAGE_REL_BASED_HIGHLOW        RelocationType = 3
	IMAGE_REL_BASED_HIGHADJ        RelocationType = 4
	IMAGE_REL_BASED_MIPS_JMPADDR   RelocationType = 5
	IMAGE_REL_BASED_SECTION        RelocationType = 6
	IMAGE_REL_BASED_REL            RelocationType = 7
	IMAGE_REL_BASED_MIPS_JMPADDR16 RelocationType = 9
	IMAGE_REL_BASED_IA64_IMM64     RelocationType = 9
	IMAGE_REL_BASED_DIR64          RelocationType = 10
	IMAGE_REL_BASED_HIGH3ADJ       RelocationType = 11
)

type BaseRelocation struct {
	VirtualAddress uint32
	SizeOfBlock    uint32
}

type BaseRelocationEntry struct {
	Data uint16
}

type TlsDirectory struct {
	StartAddressOfRawData uint32
	EndAddressOfRawData   uint32
	AddressOfIndex        uint32
	AddressOfCallBacks    uint32
	SizeOfZeroFill        uint32
	Characteristics       TlsDirectoryCharacteristics
}

type TlsDirectory64 struct {
	StartAddressOfRawData uint64
	EndAddressOfRawData   uint64
	AddressOfIndex        uint64
	AddressOfCallBacks    uint64
	SizeOfZeroFill        uint32
	Characteristics       TlsDirectoryCharacteristics
}

type LoadConfigDirectory struct {
	Size                          uint32
	TimeDateStamp                 TimeDateStamp
	MajorVersion                  uint16
	MinorVersion                  uint16
	GlobalFlagsClear              uint32
	GlobalFlagsSet                uint32
	CriticalSectionDefaultTimeout uint32
	DeCommitFreeBlockThreshold    uint32
	DeCommitTotalFreeThreshold    uint32
	LockPrefixTable               uint32
	MaximumAllocationSize         uint32
	VirtualMemoryThreshold        uint32
	ProcessHeapFlags              uint32
	ProcessAffinityMask           uint32
	CSDVersion                    uint16
	Reserved1                     uint16
	EditList                      uint32
	SecurityCookie                uint32
	SEHandlerTable                uint32
	SEHandlerCount                uint32
	GuardCFCheckFunctionPointer   uint32
	Reserved2                     uint32
	GuardCFFunctionTable          uint32
	GuardCFFunctionCount          uint32
	GuardFlags                    uint32
}

type LoadConfigDirectory64 struct {
	Size                          uint32
	TimeDateStamp                 TimeDateStamp
	MajorVersion                  uint16
	MinorVersion                  uint16
	GlobalFlagsClear              uint32
	GlobalFlagsSet                uint32
	CriticalSectionDefaultTimeout uint32
	DeCommitFreeBlockThreshold    uint64
	DeCommitTotalFreeThreshold    uint64
	LockPrefixTable               uint64
	MaximumAllocationSize         uint64
	VirtualMemoryThreshold        uint64
	ProcessAffinityMask           uint64
	ProcessHeapFlags              uint32
	CSDVersion                    uint16
	Reserved1                     uint16
	EditList                      uint64
	SecurityCookie                uint64
	SEHandlerTable                uint64
	SEHandlerCount                uint64
	GuardCFCheckFunctionPointer   uint64
	Reserved2                     uint64
	GuardCFFunctionTable          uint64
	GuardCFFunctionCount          uint64
	GuardFlags                    uint32
}

type BoundImportDescriptor struct {
	TimeDateStamp               TimeDateStamp
	OffsetModuleName            uint16
	NumberOfModuleForwarderRefs uint16
}

type BoundForwarderRef struct {
	TimeDateStamp    TimeDateStamp
	OffsetModuleName uint16
	Reserved         uint16
}

type RichHeader struct {
	Magic1    uint32 // DanS ^ checksum
	Checksums [3]uint32
	Magic2    uint32 // Rich
}
