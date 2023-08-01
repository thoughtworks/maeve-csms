cert_dir := config/certificates

# Temporary file for process substitution
temp_file := $(cert_dir)/temp.ext

csms.key:
	openssl ecparam -name prime256v1 -genkey -noout -out $(cert_dir)/csms.key

csms.csr: csms.key
	openssl req -new -nodes -key $(cert_dir)/csms.key \
		-subj "/CN=CSMS/O=Thoughtworks" \
		-addext "subjectAltName = DNS:localhost, DNS:gateway, DNS:lb" \
		-out $(cert_dir)/csms.csr

csms.pem: csms.csr
	echo "basicConstraints = critical, CA:false" > $(temp_file)
	echo "keyUsage = critical, digitalSignature, keyEncipherment" >> $(temp_file)
	echo "subjectAltName = DNS:localhost, DNS:gateway, DNS:lb" >> $(temp_file)
	openssl x509 -req -in $(cert_dir)/csms.csr \
		-out $(cert_dir)/csms.pem \
		-signkey $(cert_dir)/csms.key \
		-days 365 \
		-extfile $(temp_file)
	rm -f $(temp_file)

.PHONY: clean
clean:
	rm -f $(cert_dir)/*
