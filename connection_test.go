package libvirt

import (
	"bytes"
	"testing"
)

const (
	QEMUSystemURI  = "qemu:///system"
	TestDefaultURI = "test:///default"
)

func openTestConnection(t testing.TB) Connection {
	conn, err := Open(QEMUSystemURI)
	if err != nil {
		t.Fatal(err)
	}

	return conn
}

func TestOpen(t *testing.T) {
	if _, err := Open("xxx"); err == nil {
		t.Error("an error was not returned when connecting to a bad URI")
	}

	conn, err := Open(QEMUSystemURI)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	if !conn.IsAlive() {
		t.Error("the libvirt connection was opened but it is not alive")
	}

	// IsEncrypted

	if !conn.IsSecure() {
		t.Error("the libvirt connection is not secure")
	}
}

func TestOpenReadOnly(t *testing.T) {
	if _, err := OpenReadOnly("xxx"); err == nil {
		t.Error("an error was not returned when connecting (RO) to a bad URI")
	}

	conn, err := OpenReadOnly(QEMUSystemURI)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	if !conn.IsAlive() {
		t.Error("the libvirt connection was opened but it is not alive")
	}

	// IsEncrypted

	if !conn.IsSecure() {
		t.Error("the libvirt connection is not secure")
	}
}

func TestVersion(t *testing.T) {
	conn := openTestConnection(t)
	defer conn.Close()

	version, err := conn.Version()
	if err != nil {
		t.Fatal(err)
	}

	if version < 0 {
		t.Errorf("hypervisor version should be a positive number: %d", version)
	}
}

func TestLibVersion(t *testing.T) {
	conn := openTestConnection(t)
	defer conn.Close()

	version, err := conn.LibVersion()
	if err != nil {
		t.Fatal(err)
	}

	if version < 0 {
		t.Errorf("libvirt version should be a positive number: %d", version)
	}
}

func TestCapabilities(t *testing.T) {
	conn := openTestConnection(t)
	defer conn.Close()

	cap, err := conn.Capabilities()
	if err != nil {
		t.Fatal(err)
	}

	if len(cap) == 0 {
		t.Error("libvirt capabilities should not be empty")
	}
}

func TestHostname(t *testing.T) {
	conn := openTestConnection(t)
	defer conn.Close()

	hostname, err := conn.Hostname()
	if err != nil {
		t.Fatal(err)
	}

	if len(hostname) == 0 {
		t.Error("libvirt hostname should not be empty")
	}
}

func TestSysinfo(t *testing.T) {
	conn := openTestConnection(t)
	defer conn.Close()

	sysinfo, err := conn.Sysinfo()
	if err != nil {
		t.Fatal(err)
	}

	if len(sysinfo) == 0 {
		t.Error("libvirt sysinfo should not be empty")
	}
}

func TestType(t *testing.T) {
	conn := openTestConnection(t)
	defer conn.Close()

	typ, err := conn.Type()
	if err != nil {
		t.Fatal(err)
	}

	if len(typ) == 0 {
		t.Error("libvirt type should not be empty")
	}
}

func TestURI(t *testing.T) {
	conn := openTestConnection(t)
	defer conn.Close()

	uri, err := conn.URI()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte(uri), []byte(QEMUSystemURI)) {
		t.Errorf("libvirt URI should be the same used to open the connection; got=%s, want=%s", uri, QEMUSystemURI)
	}
}

func TestRef(t *testing.T) {
	conn := openTestConnection(t)

	if err := conn.Ref(); err != nil {
		t.Fatal(err)
	}

	if _, err := conn.Close(); err != nil {
		t.Error(err)
	}
	if _, err := conn.Close(); err != nil {
		t.Error("could not close the connection for the second time after calling Ref")
	}
}

func TestCPUModelNames(t *testing.T) {
	conn := openTestConnection(t)
	defer conn.Close()

	if _, err := conn.CPUModelNames("xxx"); err == nil {
		t.Error("an error was not returned when getting CPU model names from invalid arch")
	}

	models, err := conn.CPUModelNames("x86_64")
	if err != nil {
		t.Fatal(err)
	}

	if len(models) == 0 {
		t.Error("libvirt CPU model names should not be empty")
	}
}

func TestMaxVCPUs(t *testing.T) {
	conn := openTestConnection(t)
	defer conn.Close()

	if _, err := conn.MaxVCPUs("xxx"); err == nil {
		t.Error("an error was not returned when getting maximum VCPUs from invalid type")
	}

	vcpus, err := conn.MaxVCPUs("kvm")
	if err != nil {
		t.Fatal(err)
	}

	if vcpus < 0 {
		t.Error("libvirt maximum VCPU should be a positive number")
	}
}

func TestListDomains(t *testing.T) {
	conn := openTestConnection(t)
	defer conn.Close()

	domains, err := conn.ListDomains(DomAll)
	if err != nil {
		t.Fatal(err)
	}

	for _, d := range domains {
		if err := d.Free(); err != nil {
			t.Error(err)
		}
	}
}

func BenchmarkConnection(b *testing.B) {
	for n := 0; n < b.N; n++ {
		conn, err := Open(QEMUSystemURI)
		if err != nil {
			b.Error(err)
		}

		if _, err := conn.Close(); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkTestConnection(b *testing.B) {
	for n := 0; n < b.N; n++ {
		conn, err := Open(TestDefaultURI)
		if err != nil {
			b.Error(err)
		}

		if _, err := conn.Close(); err != nil {
			b.Error(err)
		}
	}
}