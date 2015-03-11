package intercom

import "time"
import "fmt"

/*
user_id	yes if no email	a unique string identifier for the user. It is required on creation if an email is not supplied.
email	yes if no user_id	the user’s email address. It is required on creation if a user_id is not supplied.
id	no	The id may be used for user updates.
signed_up_at	timestamp	The time the user signed up
name	no	The user’s full name
last_seen_ip	no	An ip address (e.g. “1.2.3.4”) representing the last ip address the user visited your application from. (Used for updating location_data)
custom_attributes	no	A hash of key/value pairs containing any other data about the user you want Intercom to store.*
last_seen_user_agent	no	The user agent the user last visited your application with.
companies	no	Identifies the companies this user belongs to.
last_request_at	no	A UNIX timestamp representing the date the user last visited your application.
unsubscribed_from_emails	no	A boolean value representing the users unsubscribed status. default value if not sent is false.
update_last_request_at	no	A boolean value, which if true, instructs Intercom to update the users' last_request_at value to the current API service time in UTC. default value if not sent is false.
new_session	no	A boolean value, which if true, instructs Intercom to register the request as a session.
*/
type UserUpdate struct {
	UserId              string
	Email               string
	SignedUpAt          time.Time
	Name                string
	LastSeenIp          string
	CustomAttributes    map[string]interface{}
	LastSeenUserAgent   string
	LastRequestAt       time.Time
	Unsubscribed        bool
	UpdateLastRequestAt bool
	NewSession          bool
}

type RateLimitError int64

func (i RateLimitError) Error() string {
	return fmt.Sprintf("Rate limit will reset at: '%d'", i)
}
