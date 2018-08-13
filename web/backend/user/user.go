package user

type User interface {
    GetUserId() string
    GetUsername() string
}

type DatabaseUser struct {
    ServiceName string      // For now, only "twitter"
    ServiceId string        // Identifier from the service
    UserId string           // Unique id we come up with, a UUID
    Username string
}

func (d DatabaseUser) GetUserId() string {
    return d.UserId
}

func (d DatabaseUser) GetUsername() string {
    return d.Username
}

