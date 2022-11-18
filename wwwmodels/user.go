package wwwmodels

import "database/sql"

type User struct {
	Userid int 
	User_hearabout_id int 
	User_leadsource_id int 
	User_zones int 
	Usertype_id int 
	Znode_id int
	Email, Preferred_lang, Create_dt, Last_mod_dt string
	Password, Fname, Lname, Ship_address1, Ship_address2, Ship_city, Ship_state, Ship_postal_code, Ship_country sql.NullString
	Home_phone_cc, Home_phone_num, Office_phone_cc, Office_phone_num, Office_phone_ext, Mobile_phone_cc, Mobile_phone_num sql.NullString
	Preferred_locale, Channelwave_id, Gp_id, Create_by, Last_mod_by, SoftwareVerMax, PasswordSalt, LastLoginDt, PwdLastModDt sql.NullString
	Last_mod_IP sql.NullString
	Rightnow_id, Storefront_id, Forums_id, Email_verification_grace_period, Okta_userid, Okta_first_login_dt, SalesforceAccount_id, IsCustomerDate sql.NullInt64
	Email_critical_updates, Email_keep_informed, Email_digital_expert, IsCustomer, Rightnow_update, Channelwave_update, Storefront_update bool
	Forums_update, IsBeta, IsDE, Ship_country_setbysoft, Email_survey, IsEncrypted, Znode_update bool
	PwdIsUserSet, IsEmailVerified, Inappmessaging sql.NullBool
}