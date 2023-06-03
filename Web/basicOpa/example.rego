package example.authz

import future.keywords.if
import future.keywords.in

default allow := false

allow if {
    input.method == "GET"
    input.path == ["salary", input.subject.user]
    input.expiry_year >= 2020
}

is_admin if input.subject.user == "admin"

allow if is_admin
