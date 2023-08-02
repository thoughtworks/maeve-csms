# Temporary file for process substitution
temp_file := /tmp/temp.ext

all: csms.pem cpo_sub_ca1.pem cpo_sub_ca2.pem root-V2G-cert.pem trust.pem

csms.key:
	openssl ecparam -name prime256v1 -genkey -noout -out csms.key

csms.csr: csms.key
	openssl req -new -nodes -key csms.key \
		-subj "/CN=CSMS/O=Thoughtworks" \
		-addext "subjectAltName = DNS:localhost, DNS:gateway, DNS:lb" \
		-out csms.csr

csms.pem: csms.csr
	echo "basicConstraints = critical, CA:false" > $(temp_file)
	echo "keyUsage = critical, digitalSignature, keyEncipherment" >> $(temp_file)
	echo "subjectAltName = DNS:localhost, DNS:gateway, DNS:lb" >> $(temp_file)
	openssl x509 -req -in csms.csr \
		-out csms.pem \
		-signkey csms.key \
		-days 365 \
		-extfile $(temp_file)
	rm -f $(temp_file)

cpo_sub_ca1.pem cpo_sub_ca2.pem root-V2G-cert.pem trust.pem:
	../../scripts/get-ca-cert.sh

.PHONY: clean
clean:
	rm -f *.pem csms.key csms.csr
