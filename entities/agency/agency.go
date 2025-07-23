package agency

import (
	"agency/entities/region"
	"agency/utils"
	"fmt"
	"strconv"
)

type Agency struct {
	UUID, Name, Address, PhoneNumber, MembershipDate string
	NumberOfEmployees                                uint64
	region.Region
}

func New(name, address, phoneNumber, membershipDate, regionName string, numberOfEmployees uint64) *Agency {
	return &Agency{
		UUID:              utils.GenerateUUID(),
		Name:              name,
		Address:           address,
		PhoneNumber:       phoneNumber,
		MembershipDate:    membershipDate,
		NumberOfEmployees: numberOfEmployees,
		Region: region.Region{
			RegionName: regionName,
		},
	}
}

func (a Agency) String() string {
	const formatText = "%s Agency in the %s Region\nAddress: %s\nPhone Number: %s\nMembership Date: %s\nEmployees: %d\n"
	return fmt.Sprintf(formatText, a.Name, a.RegionName, a.Address, a.PhoneNumber, a.MembershipDate, a.NumberOfEmployees)
}

func ToCSV(a *Agency) []string {
	return []string{
		a.UUID,
		a.Name,
		a.Address,
		a.PhoneNumber,
		a.MembershipDate,
		strconv.FormatUint(a.NumberOfEmployees, 10),
		a.RegionName,
	}
}

func FromCSV(record []string) *Agency {
	NumberOfEmployees, _ := strconv.ParseUint(record[5], 10, 64)
	return &Agency{
		UUID:              record[0],
		Name:              record[1],
		Address:           record[2],
		PhoneNumber:       record[3],
		MembershipDate:    record[4],
		NumberOfEmployees: NumberOfEmployees,
		Region: region.Region{
			RegionName: record[6],
		},
	}
}
