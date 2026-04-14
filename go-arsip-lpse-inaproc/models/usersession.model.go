package models


type UserSession struct {
	Id 		uint
	Name 	string
	Role 	string
}

func (c UserSession) IsPpk() bool {
	return len(c.Role) > 0 && c.Role == PPK
}

func (c UserSession) IsPokja() bool {
	return len(c.Role) > 0 && c.Role == POKJA
}

func (c UserSession) IsPp() bool {
	return len(c.Role) > 0 && c.Role == PP
}

func (c UserSession) IsAdmin() bool {
	return len(c.Role) > 0 && c.Role == ADMIN
}


func (c UserSession) IsAdmAgency() bool {
	return len(c.Role) > 0 && c.Role == ADMIN
}

func (c UserSession) IsUkpbj() bool {
	return len(c.Role) > 0 && (c.Role == UKPBJ || c.Role == ARSIPARIS)
}

func (c UserSession) IsArsiparis() bool {
	return len(c.Role) > 0 && c.Role == ARSIPARIS
}

func (c UserSession) Pegawai() Pegawai {
	var user Pegawai
	db.First(&user, c.Id)
	return user
}
