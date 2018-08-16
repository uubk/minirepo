package minirepo

import "testing"

func TestMetaDownload(t *testing.T) {
	uut := NewRepoClient("/tmp", "http://127.0.0.1:8080", `-----BEGIN PGP PUBLIC KEY BLOCK-----

xsBNBFt1Wl8BCADFae0S0pdNG1EtYt4+NUyhZtLRsy6AMPGxfbalVjIYJ/w7k0hl
hXEs8r7f6/rcSeF0kFn1aOYsQK8KWzf/2XUnXozPsDZV4ohV0ke4TvKjMuFTIzEb
BX8vs5Nxuq54GJSN0hNPBxMSlcuxm38AbBA9tr2qhSNbB9310+34sXIbcqN48DUt
8811Htsq5CI/mSvHhfEd2Yp9hLkOtLGl225sPRSoVhCoa8WUSrffRrGVPjDUCNeH
lV+wj7Vic2fTSWJEhRIE3WzvG0UrwLT9xXQIpDLJLomr7dZZ1UTvyZO80Y3J+fFs
Dk3WKaitVlow1Dk06ifFnF5QUmuGRQG4zRXDABEBAAHNGG1pbmlyZXBvIChBdXRv
Z2VuZXJhdGVkKcLAYgQTAQkAFgUCW3VaXwkQwX/uVu8e1nQCGwMCGQEAAPOoCAAv
RcRatNjmPal8gFqbxMo0rRIa6zqMpyGkUekfW6MQIVWCmH0uMjkq6p+lm8CG91+T
/NhtpDZD4uh5JSS5YwR/M1QKd6YbwIc5ZbbEZoRwJ3VS08AHm1xTt6IXSdj0nxnX
J/BxmZ0WJLu3ORbk+fsxjCVbXstB/eJtFzgxZNKrtWWE8GktxeJHWhyXhNbmKwAW
us2BerAuATAlSk2M/WAW1s8/OuAasXh13PRVeMPc8NVTpZxxh9zJuGpp+Tf9lwwS
PtwCPRalK4II+/aODxkP+FiPb2mDTbegDz+tVY8mLhRQ3KbSzTv9Sh0FEdufLF0j
2U0NJnOFRfYjdBQ1Vby7zsBNBFt1Wl8BCAC8fO6x3LC5Jo15drM2IlCKUKnX7Hrp
GiuXZeEL8G0R5f87tbvXKCJ2e+ZveaDC1XfFGGC9EkuIpNIzp2kQ6sX3kZ1etOt6
AU+Aj/SAp78b5OHXMyr4W94HXj7pjJ1N3uEoy/Z52HBYm+TJ60BIor3AJp2L4Q0g
pwlasusjvCbxSMSRBq3cxpyUv3B8EONrN9uIJ8lg5Gm/11/BzbfxMDYbbNi3qnm9
I4PBzpk6RHiSFgkHezkWlYapoSfzlceK4/ArFkDCgJH47AZaeRcS3CAu+vLCAbkS
luZ32uN0dQsC79eu+he2Nr1Q7Ik3ZBg0Z4kXY/QxxkFU9oAOsrH6Ip8tABEBAAHC
wF8EGAEJABMFAlt1Wl8JEMF/7lbvHtZ0AhsMAABjjQgAM6av5av0WCQqNc/lFhRL
SdrmyuQOWAkY00/tomvsBpvpuQl9mvge77mrzKWAkIWNxFMBdssZ7/BTzBlBQt7/
bqvDdI9tbMb2q2o86oLbjEHaNb4wLF6QTaif2eovpA/eHUFj/3GO+RwUWN2wYrRz
T584BMKEqIKxk4lBACzYVCoapzWnm/cULUGtfzps0TdyoqieWpaBFbw0/6Bt0rVq
hh1J4XZBoQuv9DbRsBWdExx1kISddmmNwTqm7ZZhyZdhdSVCWYfLlORiG6e0g3ro
cXON3imVgTBW7VkbDic6B8XVDvVY5pmkvB/HNtkTfGyb6EnaMCscfO3ipnNhXa4Q
FQ==
=Le4N
-----END PGP PUBLIC KEY BLOCK-----`)
	err := uut.fetchMeta()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
}

func TestFileDownload(t *testing.T) {
	uut := NewRepoClient("/tmp", "http://127.0.0.1:8080", `-----BEGIN PGP PUBLIC KEY BLOCK-----

xsBNBFt1Wl8BCADFae0S0pdNG1EtYt4+NUyhZtLRsy6AMPGxfbalVjIYJ/w7k0hl
hXEs8r7f6/rcSeF0kFn1aOYsQK8KWzf/2XUnXozPsDZV4ohV0ke4TvKjMuFTIzEb
BX8vs5Nxuq54GJSN0hNPBxMSlcuxm38AbBA9tr2qhSNbB9310+34sXIbcqN48DUt
8811Htsq5CI/mSvHhfEd2Yp9hLkOtLGl225sPRSoVhCoa8WUSrffRrGVPjDUCNeH
lV+wj7Vic2fTSWJEhRIE3WzvG0UrwLT9xXQIpDLJLomr7dZZ1UTvyZO80Y3J+fFs
Dk3WKaitVlow1Dk06ifFnF5QUmuGRQG4zRXDABEBAAHNGG1pbmlyZXBvIChBdXRv
Z2VuZXJhdGVkKcLAYgQTAQkAFgUCW3VaXwkQwX/uVu8e1nQCGwMCGQEAAPOoCAAv
RcRatNjmPal8gFqbxMo0rRIa6zqMpyGkUekfW6MQIVWCmH0uMjkq6p+lm8CG91+T
/NhtpDZD4uh5JSS5YwR/M1QKd6YbwIc5ZbbEZoRwJ3VS08AHm1xTt6IXSdj0nxnX
J/BxmZ0WJLu3ORbk+fsxjCVbXstB/eJtFzgxZNKrtWWE8GktxeJHWhyXhNbmKwAW
us2BerAuATAlSk2M/WAW1s8/OuAasXh13PRVeMPc8NVTpZxxh9zJuGpp+Tf9lwwS
PtwCPRalK4II+/aODxkP+FiPb2mDTbegDz+tVY8mLhRQ3KbSzTv9Sh0FEdufLF0j
2U0NJnOFRfYjdBQ1Vby7zsBNBFt1Wl8BCAC8fO6x3LC5Jo15drM2IlCKUKnX7Hrp
GiuXZeEL8G0R5f87tbvXKCJ2e+ZveaDC1XfFGGC9EkuIpNIzp2kQ6sX3kZ1etOt6
AU+Aj/SAp78b5OHXMyr4W94HXj7pjJ1N3uEoy/Z52HBYm+TJ60BIor3AJp2L4Q0g
pwlasusjvCbxSMSRBq3cxpyUv3B8EONrN9uIJ8lg5Gm/11/BzbfxMDYbbNi3qnm9
I4PBzpk6RHiSFgkHezkWlYapoSfzlceK4/ArFkDCgJH47AZaeRcS3CAu+vLCAbkS
luZ32uN0dQsC79eu+he2Nr1Q7Ik3ZBg0Z4kXY/QxxkFU9oAOsrH6Ip8tABEBAAHC
wF8EGAEJABMFAlt1Wl8JEMF/7lbvHtZ0AhsMAABjjQgAM6av5av0WCQqNc/lFhRL
SdrmyuQOWAkY00/tomvsBpvpuQl9mvge77mrzKWAkIWNxFMBdssZ7/BTzBlBQt7/
bqvDdI9tbMb2q2o86oLbjEHaNb4wLF6QTaif2eovpA/eHUFj/3GO+RwUWN2wYrRz
T584BMKEqIKxk4lBACzYVCoapzWnm/cULUGtfzps0TdyoqieWpaBFbw0/6Bt0rVq
hh1J4XZBoQuv9DbRsBWdExx1kISddmmNwTqm7ZZhyZdhdSVCWYfLlORiG6e0g3ro
cXON3imVgTBW7VkbDic6B8XVDvVY5pmkvB/HNtkTfGyb6EnaMCscfO3ipnNhXa4Q
FQ==
=Le4N
-----END PGP PUBLIC KEY BLOCK-----`)
	err := uut.fetchMeta()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	err = uut.decodeMeta()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	file, err := uut.GetFile("foo", "bar", "test")
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	print(file)
}
