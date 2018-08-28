package pefile

import (
	"time"
)

const (
	IMAGE_DOS_SIGNATURE   = 0x5A4D
	IMAGE_DOSZM_SIGNATURE = 0x4D5A
	IMAGE_NE_SIGNATURE    = 0x454E
	IMAGE_LE_SIGNATURE    = 0x454C
	IMAGE_LX_SIGNATURE    = 0x584C
	IMAGE_TE_SIGNATURE    = 0x5A56 // Terse Executables have a 'VZ' signature
	IMAGE_NT_SIGNATURE    = 0x00004550
)

type OptionalHeaderMagic uint16

const (
	OPTIONAL_HEADER_MAGIC_PE      OptionalHeaderMagic = 0x10b
	OPTIONAL_HEADER_MAGIC_PE_PLUS OptionalHeaderMagic = 0x20b
)

const (
	// IMAGE_RICH_KEY = A3 D2 F3 B4
	IMAGE_RICH_KEY       = 0x92033d19
	IMAGE_RICH_SIGNATURE = 0x68636952
	IMAGE_DANS_SIGNATURE = 0x68636952
)

type TimeDateStamp uint32

func (ts TimeDateStamp) String() string {
	ut := time.Unix(int64(ts), 0)
	return ut.Format(time.RFC3339)
}

type SubsystemType uint16

const (
	IMAGE_SUBSYSTEM_UNKNOWN                  SubsystemType = 0
	IMAGE_SUBSYSTEM_NATIVE                   SubsystemType = 1
	IMAGE_SUBSYSTEM_WINDOWS_GUI              SubsystemType = 2
	IMAGE_SUBSYSTEM_WINDOWS_CUI              SubsystemType = 3
	IMAGE_SUBSYSTEM_OS2_CUI                  SubsystemType = 5
	IMAGE_SUBSYSTEM_POSIX_CUI                SubsystemType = 7
	IMAGE_SUBSYSTEM_NATIVE_WINDOWS           SubsystemType = 8
	IMAGE_SUBSYSTEM_WINDOWS_CE_GUI           SubsystemType = 9
	IMAGE_SUBSYSTEM_EFI_APPLICATION          SubsystemType = 10
	IMAGE_SUBSYSTEM_EFI_BOOT_SERVICE_DRIVER  SubsystemType = 11
	IMAGE_SUBSYSTEM_EFI_RUNTIME_DRIVER       SubsystemType = 12
	IMAGE_SUBSYSTEM_EFI_ROM                  SubsystemType = 13
	IMAGE_SUBSYSTEM_XBOX                     SubsystemType = 14
	IMAGE_SUBSYSTEM_WINDOWS_BOOT_APPLICATION SubsystemType = 16
)

type MachineType uint16

const (
	IMAGE_FILE_MACHINE_UNKNOWN   MachineType = 0
	IMAGE_FILE_MACHINE_I386      MachineType = 0x014c
	IMAGE_FILE_MACHINE_R3000     MachineType = 0x0162
	IMAGE_FILE_MACHINE_R4000     MachineType = 0x0166
	IMAGE_FILE_MACHINE_R10000    MachineType = 0x0168
	IMAGE_FILE_MACHINE_WCEMIPSV2 MachineType = 0x0169
	IMAGE_FILE_MACHINE_ALPHA     MachineType = 0x0184
	IMAGE_FILE_MACHINE_SH3       MachineType = 0x01a2
	IMAGE_FILE_MACHINE_SH3DSP    MachineType = 0x01a3
	IMAGE_FILE_MACHINE_SH3E      MachineType = 0x01a4
	IMAGE_FILE_MACHINE_SH4       MachineType = 0x01a6
	IMAGE_FILE_MACHINE_SH5       MachineType = 0x01a8
	IMAGE_FILE_MACHINE_ARM       MachineType = 0x01c0
	IMAGE_FILE_MACHINE_THUMB     MachineType = 0x01c2
	IMAGE_FILE_MACHINE_ARMNT     MachineType = 0x01c4
	IMAGE_FILE_MACHINE_AM33      MachineType = 0x01d3
	IMAGE_FILE_MACHINE_POWERPC   MachineType = 0x01f0
	IMAGE_FILE_MACHINE_POWERPCFP MachineType = 0x01f1
	IMAGE_FILE_MACHINE_IA64      MachineType = 0x0200
	IMAGE_FILE_MACHINE_MIPS16    MachineType = 0x0266
	IMAGE_FILE_MACHINE_ALPHA64   MachineType = 0x0284
	IMAGE_FILE_MACHINE_AXP64     MachineType = 0x0284 // same
	IMAGE_FILE_MACHINE_MIPSFPU   MachineType = 0x0366
	IMAGE_FILE_MACHINE_MIPSFPU16 MachineType = 0x0466
	IMAGE_FILE_MACHINE_TRICORE   MachineType = 0x0520
	IMAGE_FILE_MACHINE_CEF       MachineType = 0x0cef
	IMAGE_FILE_MACHINE_EBC       MachineType = 0x0ebc
	IMAGE_FILE_MACHINE_AMD64     MachineType = 0x8664
	IMAGE_FILE_MACHINE_M32R      MachineType = 0x9041
	IMAGE_FILE_MACHINE_CEE       MachineType = 0xc0ee
)

type ImageCharacteristics uint16

const (
	IMAGE_FILE_RELOCS_STRIPPED         ImageCharacteristics = 0x0001
	IMAGE_FILE_EXECUTABLE_IMAGE        ImageCharacteristics = 0x0002
	IMAGE_FILE_LINE_NUMS_STRIPPED      ImageCharacteristics = 0x0004
	IMAGE_FILE_LOCAL_SYMS_STRIPPED     ImageCharacteristics = 0x0008
	IMAGE_FILE_AGGRESIVE_WS_TRIM       ImageCharacteristics = 0x0010
	IMAGE_FILE_LARGE_ADDRESS_AWARE     ImageCharacteristics = 0x0020
	IMAGE_FILE_16BIT_MACHINE           ImageCharacteristics = 0x0040
	IMAGE_FILE_BYTES_REVERSED_LO       ImageCharacteristics = 0x0080
	IMAGE_FILE_32BIT_MACHINE           ImageCharacteristics = 0x0100
	IMAGE_FILE_DEBUG_STRIPPED          ImageCharacteristics = 0x0200
	IMAGE_FILE_REMOVABLE_RUN_FROM_SWAP ImageCharacteristics = 0x0400
	IMAGE_FILE_NET_RUN_FROM_SWAP       ImageCharacteristics = 0x0800
	IMAGE_FILE_SYSTEM                  ImageCharacteristics = 0x1000
	IMAGE_FILE_DLL                     ImageCharacteristics = 0x2000
	IMAGE_FILE_UP_SYSTEM_ONLY          ImageCharacteristics = 0x4000
	IMAGE_FILE_BYTES_REVERSED_HI       ImageCharacteristics = 0x8000
)

type DllCharacteristics uint16

const (
	IMAGE_LIBRARY_PROCESS_INIT                     DllCharacteristics = 0x0001 // reserved
	IMAGE_LIBRARY_PROCESS_TERM                     DllCharacteristics = 0x0002 // reserved
	IMAGE_LIBRARY_THREAD_INIT                      DllCharacteristics = 0x0004 // reserved
	IMAGE_LIBRARY_THREAD_TERM                      DllCharacteristics = 0x0008 // reserved
	IMAGE_DLLCHARACTERISTICS_HIGH_ENTROPY_VA       DllCharacteristics = 0x0020
	IMAGE_DLLCHARACTERISTICS_DYNAMIC_BASE          DllCharacteristics = 0x0040
	IMAGE_DLLCHARACTERISTICS_FORCE_INTEGRITY       DllCharacteristics = 0x0080
	IMAGE_DLLCHARACTERISTICS_NX_COMPAT             DllCharacteristics = 0x0100
	IMAGE_DLLCHARACTERISTICS_NO_ISOLATION          DllCharacteristics = 0x0200
	IMAGE_DLLCHARACTERISTICS_NO_SEH                DllCharacteristics = 0x0400
	IMAGE_DLLCHARACTERISTICS_NO_BIND               DllCharacteristics = 0x0800
	IMAGE_DLLCHARACTERISTICS_APPCONTAINER          DllCharacteristics = 0x1000
	IMAGE_DLLCHARACTERISTICS_WDM_DRIVER            DllCharacteristics = 0x2000
	IMAGE_DLLCHARACTERISTICS_GUARD_CF              DllCharacteristics = 0x4000
	IMAGE_DLLCHARACTERISTICS_TERMINAL_SERVER_AWARE DllCharacteristics = 0x8000
)

type SectionCharacteristics uint32

const (
	IMAGE_SCN_TYPE_REG               SectionCharacteristics = 0x00000000 // reserved
	IMAGE_SCN_TYPE_DSECT             SectionCharacteristics = 0x00000001 // reserved
	IMAGE_SCN_TYPE_NOLOAD            SectionCharacteristics = 0x00000002 // reserved
	IMAGE_SCN_TYPE_GROUP             SectionCharacteristics = 0x00000004 // reserved
	IMAGE_SCN_TYPE_NO_PAD            SectionCharacteristics = 0x00000008 // reserved
	IMAGE_SCN_TYPE_COPY              SectionCharacteristics = 0x00000010 // reserved
	IMAGE_SCN_CNT_CODE               SectionCharacteristics = 0x00000020
	IMAGE_SCN_CNT_INITIALIZED_DATA   SectionCharacteristics = 0x00000040
	IMAGE_SCN_CNT_UNINITIALIZED_DATA SectionCharacteristics = 0x00000080
	IMAGE_SCN_LNK_OTHER              SectionCharacteristics = 0x00000100
	IMAGE_SCN_LNK_INFO               SectionCharacteristics = 0x00000200
	IMAGE_SCN_LNK_OVER               SectionCharacteristics = 0x00000400 // reserved
	IMAGE_SCN_LNK_REMOVE             SectionCharacteristics = 0x00000800
	IMAGE_SCN_LNK_COMDAT             SectionCharacteristics = 0x00001000
	IMAGE_SCN_MEM_PROTECTED          SectionCharacteristics = 0x00004000 // obsolete
	IMAGE_SCN_NO_DEFER_SPEC_EXC      SectionCharacteristics = 0x00004000
	IMAGE_SCN_GPREL                  SectionCharacteristics = 0x00008000
	IMAGE_SCN_MEM_FARDATA            SectionCharacteristics = 0x00008000
	IMAGE_SCN_MEM_SYSHEAP            SectionCharacteristics = 0x00010000 // obsolete
	IMAGE_SCN_MEM_PURGEABLE          SectionCharacteristics = 0x00020000
	IMAGE_SCN_MEM_16BIT              SectionCharacteristics = 0x00020000
	IMAGE_SCN_MEM_LOCKED             SectionCharacteristics = 0x00040000
	IMAGE_SCN_MEM_PRELOAD            SectionCharacteristics = 0x00080000
	IMAGE_SCN_ALIGN_1BYTES           SectionCharacteristics = 0x00100000
	IMAGE_SCN_ALIGN_2BYTES           SectionCharacteristics = 0x00200000
	IMAGE_SCN_ALIGN_4BYTES           SectionCharacteristics = 0x00300000
	IMAGE_SCN_ALIGN_8BYTES           SectionCharacteristics = 0x00400000
	IMAGE_SCN_ALIGN_16BYTES          SectionCharacteristics = 0x00500000 // default alignment
	IMAGE_SCN_ALIGN_32BYTES          SectionCharacteristics = 0x00600000
	IMAGE_SCN_ALIGN_64BYTES          SectionCharacteristics = 0x00700000
	IMAGE_SCN_ALIGN_128BYTES         SectionCharacteristics = 0x00800000
	IMAGE_SCN_ALIGN_256BYTES         SectionCharacteristics = 0x00900000
	IMAGE_SCN_ALIGN_512BYTES         SectionCharacteristics = 0x00A00000
	IMAGE_SCN_ALIGN_1024BYTES        SectionCharacteristics = 0x00B00000
	IMAGE_SCN_ALIGN_2048BYTES        SectionCharacteristics = 0x00C00000
	IMAGE_SCN_ALIGN_4096BYTES        SectionCharacteristics = 0x00D00000
	IMAGE_SCN_ALIGN_8192BYTES        SectionCharacteristics = 0x00E00000
	IMAGE_SCN_ALIGN_MASK             SectionCharacteristics = 0x00F00000
	IMAGE_SCN_LNK_NRELOC_OVFL        SectionCharacteristics = 0x01000000
	IMAGE_SCN_MEM_DISCARDABLE        SectionCharacteristics = 0x02000000
	IMAGE_SCN_MEM_NOT_CACHED         SectionCharacteristics = 0x04000000
	IMAGE_SCN_MEM_NOT_PAGED          SectionCharacteristics = 0x08000000
	IMAGE_SCN_MEM_SHARED             SectionCharacteristics = 0x10000000
	IMAGE_SCN_MEM_EXECUTE            SectionCharacteristics = 0x20000000
	IMAGE_SCN_MEM_READ               SectionCharacteristics = 0x40000000
	IMAGE_SCN_MEM_WRITE              SectionCharacteristics = 0x80000000
)

type ImportDescriptorCharacteristics uint32

type ExportCharacteristics uint32

type ResourceDirecotryCharacteristics uint32

type DebugDirectoryCharacteristics uint32

type TlsDirectoryCharacteristics uint32

type DebugType uint32

const (
	IMAGE_DEBUG_TYPE_UNKNOWN       DebugType = 0
	IMAGE_DEBUG_TYPE_COFF          DebugType = 1
	IMAGE_DEBUG_TYPE_CODEVIEW      DebugType = 2
	IMAGE_DEBUG_TYPE_FPO           DebugType = 3
	IMAGE_DEBUG_TYPE_MISC          DebugType = 4
	IMAGE_DEBUG_TYPE_EXCEPTION     DebugType = 5
	IMAGE_DEBUG_TYPE_FIXUP         DebugType = 6
	IMAGE_DEBUG_TYPE_OMAP_TO_SRC   DebugType = 7
	IMAGE_DEBUG_TYPE_OMAP_FROM_SRC DebugType = 8
	IMAGE_DEBUG_TYPE_BORLAND       DebugType = 9
	IMAGE_DEBUG_TYPE_RESERVED10    DebugType = 10
	IMAGE_DEBUG_TYPE_CLSID         DebugType = 11
)
