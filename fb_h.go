// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package framebuffer

// <linux/fb.h>

import "unsafe"

const (
	SyncHorHighAct  = 1 << iota // horizontal sync high active
	SyncVertHighAct             // vertical sync high active
	SyncExt                     // external sync
	SyncCompHighAct             // composite sync high active
	SyncBroadcast               // broadcast video timings
	SyncOnGreen                 // sync on green
)

const (
	VModeNonInterlaced = 0 // non interlaced
	VModeInterlaced    = 1 // interlaced
	VModeDouble        = 2 // double scan
	VModeOddFieldFirst = 4 // interlaced: top line first
	VModeMask          = 255
	VModeYWrap         = 256 // ywrap instead of panning
	VModeSmoothXPan    = 512 // smooth xpan possible (internally used)
	VModeCanUpdate     = 512 // don't update x/yoffset
)

var (
	_IO_CURSOR       int
	_IOGET_VBLANK    int
	_IO_WAITFORVSYNC int
)

func init() {
	var c fb_cursor
	var v fb_vblank

	_IO_CURSOR = _IOWR('F', 0x08, int(unsafe.Sizeof(c)))
	_IOGET_VBLANK = _IOR('F', 0x12, int(unsafe.Sizeof(v)))
	_IO_WAITFORVSYNC = _IOW('F', 0x20, 4) // 4 = sizeof(uint32)
}

type fb_fix_screeninfo struct {
	id           [16]byte  // Identification string (e.g.: "TT Builtin")
	smemstart    uint64    // Physical start address of framebuffer memory.
	smemlen      uint32    // Length of framebuffer memory.
	typ          uint32    // See __TYPE_XXX values.
	type_aux     uint32    // Interleave for interleaved planes.
	visual       uint32    // See __VISUAL_XXX values.
	xpanstep     uint16    // Zero if no hardware panning.
	ypanstep     uint16    // Zero if no hardware panning.
	ywrapstep    uint16    // Zero if no hardware ywrap.
	line_length  uint32    // Length of a line in bytes.
	mmio_start   uint64    // Physical start address of mmap'd _IO.
	mmio_len     uint32    // Length of mmap'd _IO.
	accel        uint32    // Indicate to driver which specific chip/card we have.
	capabilities uint16    // See _CAP_XXXX values.
	_            [2]uint16 // Reserved for future use.
}

// Interpretation of offset for color fields: All offsets are from the right,
// inside a "pixel" value, which is exactly 'bits_per_pixel' wide (means: you
// can use the offset as right argument to <<). A pixel afterwards is a bit
// stream and is written to video memory as-is.
//
// For pseudocolor: offset and length should be the same for all color
// components. Offset specifies the position of the least significant bit
// of the pallette index in a pixel value. Length indicates the number
// of available palette entries (i.e. # of entries = 1 << length).
type fb_bitfield struct {
	offset    uint32 // beginning of bitfield
	length    uint32 // length of bitfield
	msb_right uint32 // != 0 : Most significant bit is right
}

type fb_var_screeninfo struct {
	xres           uint32 // Visible resolution.
	yres           uint32
	xres_virtual   uint32 // Virtual resolution (viewport).
	yres_virtual   uint32
	xoffset        uint32 // Offset from virtual to visible resolution.
	yoffset        uint32
	bits_per_pixel uint32      // Bit depth.
	grayscale      uint32      // 0 = color, 1 = grayscale, >1 = FOURCC
	red            fb_bitfield // bitfield in FB mem, if true colour. Else only length is significant.
	green          fb_bitfield
	blue           fb_bitfield
	transparent    fb_bitfield
	nonstd         uint32 // non-zero = non-standard pixel format.
	activate       uint32 // See __ACTIVATE_XXXX values.
	height         uint32 // Height of picture in millimetres.
	width          uint32 // Width of picture in millimetres.
	_              uint32 // AccelFlags: obsolete

	// Timing: All values in pixclocks, except pixclock.
	// Do not change these unless you really know what you are doing.
	pixclock     uint32 // Pixel clock in picoseconds.
	left_margin  uint32 // Time from sync to picture.
	right_margin uint32 // Time from picture to sync.
	upper_margin uint32 // Time from sync to picture.
	lower_margin uint32
	hsync_len    uint32    // Length of horizontal sync.
	vsync_len    uint32    // Length of vertical sync.
	sync         uint32    // See _SYNC_XXXX values.
	vmode        uint32    // See _VMODE_XXXX values.
	rotate       uint32    // Angle of counter-clockwise rotation.
	colorspace   uint32    // Colorspace for FOURCC-based modes.
	_            [4]uint32 // Reserved for future use.
}

// Copy returns a copy of the current object.
func (v *fb_var_screeninfo) Copy() *fb_var_screeninfo {
	n := new(fb_var_screeninfo)
	n.xres = v.xres
	n.yres = v.yres
	n.xres_virtual = v.xres_virtual
	n.yres_virtual = v.yres_virtual
	n.xoffset = v.xoffset
	n.yoffset = v.yoffset
	n.bits_per_pixel = v.bits_per_pixel
	n.grayscale = v.grayscale
	n.red = v.red
	n.green = v.green
	n.blue = v.blue
	n.transparent = v.transparent
	n.nonstd = v.nonstd
	n.activate = v.activate
	n.height = v.height
	n.width = v.width
	n.pixclock = v.pixclock
	n.left_margin = v.left_margin
	n.right_margin = v.right_margin
	n.upper_margin = v.upper_margin
	n.lower_margin = v.lower_margin
	n.hsync_len = v.hsync_len
	n.vsync_len = v.vsync_len
	n.sync = v.sync
	n.vmode = v.vmode
	n.rotate = v.rotate
	n.colorspace = v.colorspace
	return n
}

type fb_cmap struct {
	start  uint32         // First entry
	len    uint32         // Number of entries
	red    unsafe.Pointer // uint16*
	green  unsafe.Pointer // uint16*
	blue   unsafe.Pointer // uint16*
	transp unsafe.Pointer // uint16* -- transparency, can be NULL
}

type fb_con2fbmap struct {
	console     uint32
	framebuffer uint32
}

type fb_vblank struct {
	flags  uint32    // vblank flags.
	count  uint32    // Counter of retraces since boot.
	vcount uint32    // Current scanline position.
	hcount uint32    // Current scandot position.
	_      [4]uint32 // Reserved for future use.
}

type fb_copyarea struct {
	dx     uint32
	dy     uint32
	width  uint32
	height uint32
	sx     uint32
	sy     uint32
}

type fb_image struct {
	dx       uint32 // Where to place image
	dy       uint32
	width    uint32 // Size of image
	height   uint32
	fg_color uint32 // Only used when a mono bitmap
	bg_color uint32
	depth    uint8          // Depth of the image
	data     unsafe.Pointer // const char* -- Pointer to image data
	cmap     fb_cmap        // color map info
}

type fbcurpos struct {
	x, y uint16
}

type fb_cursor struct {
	set    uint16         // what to set
	enable uint16         // cursor on/off
	rop    uint16         // bitop operation
	mask   unsafe.Pointer // const char* -- cursor mask bits
	hot    fbcurpos       // cursor hot spot
	image  fb_image       // Cursor image
}

const (
	_FB_ACTIVE = iota
	_FB_REL_REQ
	_FB_INACTIVE
	_FB_ACQ_REQ
)

const _MAX = 32

// _IOCTL values
const (
	_IOGET_VSCREENINFO = 0x4600
	_IOPUT_VSCREENINFO = 0x4601
	_IOGET_FSCREENINFO = 0x4602
	_IOGET_CMAP        = 0x4604
	_IOPUT_CMAP        = 0x4605
	_IOPAN_DISPLAY     = 0x4606
	_IOGET_CON2FBMAP   = 0x460F
	_IOPUT_CON2FBMAP   = 0x4610
	_IO_BLANK          = 0x4611
	_IO_ALLOC          = 0x4613
	_IO_FREE           = 0x4614
	_IOGET_GLYPH       = 0x4615
	_IOGET_HWCINFO     = 0x4616
	_IOPUT_MODEINFO    = 0x4617
	_IOGET_DISPINFO    = 0x4618
)

const (
	_TYPE_PACKED_PIXELS      = iota // Packed Pixels.
	_TYPE_PLANES                    // Non interleaved planes.
	_TYPE_INTERLEAVED_PLANES        // Interleaved planes.
	_TYPE_TEXT                      // Text/attributes.
	_TYPE_VGA_PLANES                // EGA/VGA planes.
	_TYPE_FOURCC                    // Type identified by a V4L2 FOURCC.
)

const (
	_AUX_TEXT_MDA         = 0  // Monochrome text.
	_AUX_TEXT_CGA         = 1  // CGA/EGA/VGA Color text.
	_AUX_TEXT_S3_MM_IO    = 2  // S3 MM_IO fasttext.
	_AUX_TEXT_MGA_STEP16  = 3  // MGA Millenium I: text, attr, 14 reserved bytes.
	_AUX_TEXT_MGA_STEP8   = 4  // other MGAs:      text, attr,  6 reserved bytes.
	_AUX_TEXT_SVGA_GROUP  = 8  // 8-15: SVGA tileblit compatible modes.
	_AUX_TEXT_SVGA_MASK   = 7  // lower three bits says step.
	_AUX_TEXT_SVGA_STEP2  = 8  // SVGA text mode:  text, attr.
	_AUX_TEXT_SVGA_STEP4  = 9  // SVGA text mode:  text, attr,  2 reserved bytes.
	_AUX_TEXT_SVGA_STEP8  = 10 // SVGA text mode:  text, attr,  6 reserved bytes.
	_AUX_TEXT_SVGA_STEP16 = 11 // SVGA text mode:  text, attr, 14 reserved bytes.
	_AUX_TEXT_SVGA_LAST   = 15 // reserved up to 15.
)

const (
	_AUX_VGA_PLANES_VGA4 = iota // 16 color planes (EGA/VGA).
	_AUX_VGA_PLANES_CFB4        // CFB4 in planes (VGA).
	_AUX_VGA_PLANES_CFB8        // CFB8 in planes (VGA).
)

const (
	_VISUAL_MONO01             = iota // Monochr. 1=Black 0=White.
	_VISUAL_MONO10                    // Monochr. 1=White 0=Black.
	_VISUAL_TRUECOLOR                 // True color.
	_VISUAL_PSEUDOCOLOR               // Pseudo color (like atari).
	_VISUAL_DIRECTCOLOR               // Direct color.
	_VISUAL_STATIC_PSEUDOCOLOR        // Pseudo color readonly.
	_VISUAL_FOURCC                    // Visual identified by a V4L2 FOURCC.
)

const (
	_ACCEL_NONE                = 0    // no hardware accelerator.
	_ACCEL_ATARIBLITT          = 1    // Atari Blitter
	_ACCEL_AMIGABLITT          = 2    // Amiga Blitter
	_ACCEL_S3_TR_IO64          = 3    // Cybervision64 (S3 Trio64)
	_ACCEL_NCR_77C32BLT        = 4    // RetinaZ3 (NCR 77C32BLT)
	_ACCEL_S3_VIRGE            = 5    // Cybervision64/3D (S3 ViRGE)
	_ACCEL_ATI_MACH64GX        = 6    // ATI Mach 64GX family
	_ACCEL_DEC_TGA             = 7    // DEC 21030 TGA
	_ACCEL_ATI_MACH64CT        = 8    // ATI Mach 64CT family
	_ACCEL_ATI_MACH64VT        = 9    // ATI Mach 64CT family VT class
	_ACCEL_ATI_MACH64GT        = 10   // ATI Mach 64CT family GT class
	_ACCEL_SUN_CREATOR         = 11   // Sun Creator/Creator3D
	_ACCEL_SUN_CGSIX           = 12   // Sun cg6
	_ACCEL_SUN_LEO             = 13   // Sun leo/zx
	_ACCEL_IMS_TWINTURBO       = 14   // IMS Twin Turbo
	_ACCEL_3DLABS_PERMEDIA2    = 15   // 3Dlabs Permedia 2
	_ACCEL_MATROX_MGA2064W     = 16   // Matrox MGA2064W (Millenium)
	_ACCEL_MATROX_MGA1064SG    = 17   // Matrox MGA1064SG (Mystique)
	_ACCEL_MATROX_MGA2164W     = 18   // Matrox MGA2164W (Millenium II)
	_ACCEL_MATROX_MGA2164W_AGP = 19   // Matrox MGA2164W (Millenium II)
	_ACCEL_MATROX_MGAG100      = 20   // Matrox G100 (Productiva G100)
	_ACCEL_MATROX_MGAG200      = 21   // Matrox G200 (Myst, Mill, ...)
	_ACCEL_SUN_CG14            = 22   // Sun cgfourteen
	_ACCEL_SUN_BWTWO           = 23   // Sun bwtwo
	_ACCEL_SUN_CGTHREE         = 24   // Sun cgthree
	_ACCEL_SUN_TCX             = 25   // Sun tcx
	_ACCEL_MATROX_MGAG400      = 26   // Matrox G400
	_ACCEL_NV3                 = 27   // nVidia RIVA 128
	_ACCEL_NV4                 = 28   // nVidia RIVA TNT
	_ACCEL_NV5                 = 29   // nVidia RIVA TNT2
	_ACCEL_CT_6555x            = 30   // C&T 6555x
	_ACCEL_3DFX_BANSHEE        = 31   // 3Dfx Banshee
	_ACCEL_ATI_RAGE128         = 32   // ATI Rage128 family
	_ACCEL_IGS_CYBER2000       = 33   // CyberPro 2000
	_ACCEL_IGS_CYBER2010       = 34   // CyberPro 2010
	_ACCEL_IGS_CYBER5000       = 35   // CyberPro 5000
	_ACCEL_SIS_GLAMOUR         = 36   // SiS 300/630/540
	_ACCEL_3DLABS_PERMEDIA3    = 37   // 3Dlabs Permedia 3
	_ACCEL_ATI_RADEON          = 38   // ATI Radeon family
	_ACCEL_I810                = 39   // Intel 810/815
	_ACCEL_SIS_GLAMOUR_2       = 40   // SiS 315, 650, 740
	_ACCEL_SIS_XABRE           = 41   // SiS 330 ("Xabre")
	_ACCEL_I830                = 42   // Intel 830M/845G/85x/865G
	_ACCEL_NV_10               = 43   // nVidia Arch 10
	_ACCEL_NV_20               = 44   // nVidia Arch 20
	_ACCEL_NV_30               = 45   // nVidia Arch 30
	_ACCEL_NV_40               = 46   // nVidia Arch 40
	_ACCEL_XGI_VOLARI_V        = 47   // XGI Volari V3XT, V5, V8
	_ACCEL_XGI_VOLARI_Z        = 48   // XGI Volari Z7
	_ACCEL_OMAP1610            = 49   // TI OMAP16xx
	_ACCEL_TRIDENT_TGUI        = 50   // Trident TGUI
	_ACCEL_TRIDENT_3DIMAGE     = 51   // Trident 3DImage
	_ACCEL_TRIDENT_BLADE3D     = 52   // Trident Blade3D
	_ACCEL_TRIDENT_BLADEXP     = 53   // Trident BladeXP
	_ACCEL_CIRRUS_ALPINE       = 53   // Cirrus Logic 543x/544x/5480
	_ACCEL_NEOMAGIC_NM2070     = 90   // NeoMagic NM2070
	_ACCEL_NEOMAGIC_NM2090     = 91   // NeoMagic NM2090
	_ACCEL_NEOMAGIC_NM2093     = 92   // NeoMagic NM2093
	_ACCEL_NEOMAGIC_NM2097     = 93   // NeoMagic NM2097
	_ACCEL_NEOMAGIC_NM2160     = 94   // NeoMagic NM2160
	_ACCEL_NEOMAGIC_NM2200     = 95   // NeoMagic NM2200
	_ACCEL_NEOMAGIC_NM2230     = 96   // NeoMagic NM2230
	_ACCEL_NEOMAGIC_NM2360     = 97   // NeoMagic NM2360
	_ACCEL_NEOMAGIC_NM2380     = 98   // NeoMagic NM2380
	_ACCEL_PXA3XX              = 99   // PXA3xx
	_ACCEL_SAVAGE4             = 0x80 // S3 Savage4
	_ACCEL_SAVAGE3D            = 0x81 // S3 Savage3D
	_ACCEL_SAVAGE3D_MV         = 0x82 // S3 Savage3D-MV
	_ACCEL_SAVAGE2000          = 0x83 // S3 Savage2000
	_ACCEL_SAVAGE_MX_MV        = 0x84 // S3 Savage/MX-MV
	_ACCEL_SAVAGE_MX           = 0x85 // S3 Savage/MX
	_ACCEL_SAVAGE_IX_MV        = 0x86 // S3 Savage/IX-MV
	_ACCEL_SAVAGE_IX           = 0x87 // S3 Savage/IX
	_ACCEL_PROSAVAGE_PM        = 0x88 // S3 ProSavage PM133
	_ACCEL_PROSAVAGE_KM        = 0x89 // S3 ProSavage KM133
	_ACCEL_S3TWISTER_P         = 0x8a // S3 Twister
	_ACCEL_S3TWISTER_K         = 0x8b // S3 TwisterK
	_ACCEL_SUPERSAVAGE         = 0x8c // S3 Supersavage
	_ACCEL_PROSAVAGE_DDR       = 0x8d // S3 ProSavage DDR
	_ACCEL_PROSAVAGE_DDRK      = 0x8e // S3 ProSavage DDR-K
	_ACCEL_PUV3_UNIGFX         = 0xa0 // PKUnity-v3 Unigfx
)

// Device supports FOURCC-based formats
const _CAP_FOURCC = 1

const (
	_NONSTD_HAM          = 1 // Hold-And-Modify (HAM)
	_NONSTD_REV_PIX_IN_B = 2 // order of pixels in each byte is reversed
)

const (
	_ACTIVATE_NOW      = 0   // Set values immediately (or vertical blank)
	_ACTIVATE_NXTOPEN  = 1   // Activate on next open
	_ACTIVATE_TEST     = 2   // Don't set, round up impossible
	_ACTIVATE_MASK     = 15  // Values
	_ACTIVATE_VBL      = 16  // Activate values on next vertical blank
	_CHANGE_CMAP_VBL   = 32  // Change colormap on vertical blank
	_ACTIVATE_ALL      = 64  // Change all VCs on this framebuffer
	_ACTIVATE_FORCE    = 128 // Force apply even when no change
	_ACTIVATE_INV_MODE = 256 // Invalidate videomode
)

// Display rotation support
const (
	_ROTATE_UR = iota
	_ROTATE_CW
	_ROTATE_UD
	_ROTATE_CCW
)

// const PICOS2KHZ(a) (1000000000UL/(a))
// const KHZ2PICOS(a) (1000000000UL/(a))

// VESA Blanking Levels
const (
	_VESA_NO_BLANKING = iota
	_VESA_V_SYNC_SUSPEND
	_VESA_H_SYNC_SUSPEND
	_VESA_POWERDOWN
)

const (
	_BLANK_UNBLANK        = _VESA_NO_BLANKING        // screen: unblanked, hsync: on,  vsync: on
	_BLANK_NORMAL         = _VESA_NO_BLANKING + 1    // screen: blanked,   hsync: on,  vsync: on
	_BLANK_V_SYNC_SUSPEND = _VESA_V_SYNC_SUSPEND + 1 // screen: blanked,   hsync: on,  vsync: off
	_BLANK_H_SYNC_SUSPEND = _VESA_H_SYNC_SUSPEND + 1 // screen: blanked,   hsync: off, vsync: on
	_BLANK_POWERDOWN      = _VESA_POWERDOWN + 1      // screen: blanked,   hsync: off, vsync: off
)

const (
	_VBLANK_VBLANKING   = 0x001 // currently in a vertical blank
	_VBLANK_HBLANKING   = 0x002 // currently in a horizontal blank
	_VBLANK_HAVE_VBLANK = 0x004 // vertical blanks can be detected
	_VBLANK_HAVE_HBLANK = 0x008 // horizontal blanks can be detected
	_VBLANK_HAVE_COUNT  = 0x010 // global retrace counter is available
	_VBLANK_HAVE_VCOUNT = 0x020 // the vcount field is valid
	_VBLANK_HAVE_HCOUNT = 0x040 // the hcount field is valid
	_VBLANK_VSYNCING    = 0x080 // currently in a vsync
	_VBLANK_HAVE_VSYNC  = 0x100 // verical syncs can be detected
)

// Internal HW accel.
const (
	_ROP_COPY = iota
	_ROP_XOR
)

// Hardware cursor control
const (
	_CUR_SETIMAGE = 0x01
	_CUR_SETPOS   = 0x02
	_CUR_SETHOT   = 0x04
	_CUR_SETCMAP  = 0x08
	_CUR_SETSHAPE = 0x10
	_CUR_SETSIZE  = 0x20
	_CUR_SETALL   = 0xFF
)

// Settings for the generic backlight code
const (
	_BACKLIGHT_LEVELS = 128
	_BACKLIGHT_MAX    = 0xFF
)
