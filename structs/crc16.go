package structs

import "math/bits"

type ZeroCRC16Params struct {
	Poly   uint16
	Init   uint16
	RefIn  bool
	RefOut bool
	XorOut uint16
	Name   string
}

var (
	CRC16_ARC         = ZeroCRC16Params{0x8005, 0x0000, true, true, 0x0000, "CRC-16/ARC"}
	CRC16_AUG_CCITT   = ZeroCRC16Params{0x1021, 0x1D0F, false, false, 0x0000, "CRC-16/AUG-CCITT"}
	CRC16_BUYPASS     = ZeroCRC16Params{0x8005, 0x0000, false, false, 0x0000, "CRC-16/BUYPASS"}
	CRC16_CCITT_FALSE = ZeroCRC16Params{0x1021, 0xFFFF, false, false, 0x0000, "CRC-16/CCITT-FALSE"}
	CRC16_CDMA2000    = ZeroCRC16Params{0xC867, 0xFFFF, false, false, 0x0000, "CRC-16/CDMA2000"}
	CRC16_DDS_110     = ZeroCRC16Params{0x8005, 0x800D, false, false, 0x0000, "CRC-16/DDS-110"}
	CRC16_DECT_R      = ZeroCRC16Params{0x0589, 0x0000, false, false, 0x0001, "CRC-16/DECT-R"}
	CRC16_DECT_X      = ZeroCRC16Params{0x0589, 0x0000, false, false, 0x0000, "CRC-16/DECT-X"}
	CRC16_DNP         = ZeroCRC16Params{0x3D65, 0x0000, true, true, 0xFFFF, "CRC-16/DNP"}
	CRC16_EN_13757    = ZeroCRC16Params{0x3D65, 0x0000, false, false, 0xFFFF, "CRC-16/EN-13757"}
	CRC16_GENIBUS     = ZeroCRC16Params{0x1021, 0xFFFF, false, false, 0xFFFF, "CRC-16/GENIBUS"}
	CRC16_MAXIM       = ZeroCRC16Params{0x8005, 0x0000, true, true, 0xFFFF, "CRC-16/MAXIM"}
	CRC16_MCRF4XX     = ZeroCRC16Params{0x1021, 0xFFFF, true, true, 0x0000, "CRC-16/MCRF4XX"}
	CRC16_RIELLO      = ZeroCRC16Params{0x1021, 0xB2AA, true, true, 0x0000, "CRC-16/RIELLO"}
	CRC16_T10_DIF     = ZeroCRC16Params{0x8BB7, 0x0000, false, false, 0x0000, "CRC-16/T10-DIF"}
	CRC16_TELEDISK    = ZeroCRC16Params{0xA097, 0x0000, false, false, 0x0000, "CRC-16/TELEDISK"}
	CRC16_TMS37157    = ZeroCRC16Params{0x1021, 0x89EC, true, true, 0x0000, "CRC-16/TMS37157"}
	CRC16_USB         = ZeroCRC16Params{0x8005, 0xFFFF, true, true, 0xFFFF, "CRC-16/USB"}
	CRC16_CRC_A       = ZeroCRC16Params{0x1021, 0xC6C6, true, true, 0x0000, "CRC-16/CRC-A"}
	CRC16_KERMIT      = ZeroCRC16Params{0x1021, 0x0000, true, true, 0x0000, "CRC-16/KERMIT"}
	CRC16_MODBUS      = ZeroCRC16Params{0x8005, 0xFFFF, true, true, 0x0000, "CRC-16/MODBUS"}
	CRC16_X_25        = ZeroCRC16Params{0x1021, 0xFFFF, true, true, 0xFFFF, "CRC-16/X-25"}
	CRC16_XMODEM      = ZeroCRC16Params{0x1021, 0x0000, false, false, 0x0000, "CRC-16/XMODEM"}
)

type ZeroCRC16Table struct {
	params ZeroCRC16Params
	data   []uint16
}

func NewCRC16Table(params ZeroCRC16Params) *ZeroCRC16Table {
	table := &ZeroCRC16Table{
		params: params,
		data:   make([]uint16, 256),
	}
	for n := 0; n < 256; n++ {
		crc := uint16(n) << 8
		for i := 0; i < 8; i++ {
			bit := (crc & 0x8000) != 0
			crc <<= 1
			if bit {
				crc ^= params.Poly
			}
		}
		table.data[n] = crc
	}
	return table
}

func (table *ZeroCRC16Table) initValue() uint16 {
	return uint16(table.params.Init)
}

func (table *ZeroCRC16Table) calculate(data []byte) uint16 {
	crc := table.initValue()
	for _, d := range data {
		if table.params.RefIn {
			d = bits.Reverse8(d)
		}
		crc = crc<<8 ^ table.data[byte(crc>>8)^d]
	}
	return crc
}

func (table *ZeroCRC16Table) Complete(data []byte) uint16 {
	crc := table.calculate(data)
	if table.params.RefOut {
		return bits.Reverse16(crc) ^ table.params.XorOut
	}
	return crc ^ table.params.XorOut
}
