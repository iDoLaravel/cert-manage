package ui

import (
	"crypto/x509"
	"io/ioutil"
	"os"
	"time"

	"github.com/adamdecaf/cert-manage/pkg/certutil"
	"gopkg.in/yaml.v2"
)

// observatory is a printer for the 'trust_stores_observatory'
// https://github.com/nabla-c0d3/trust_stores_observatory

var (
	observatoryFormat     = "observatory"
	observatoryTimeFormat = "2006-01-02"
)

func isObservatory(format string) bool {
	return format == observatoryFormat
}

type observatoryReport struct {
	Platform     string            `yaml:"platform"`
	Version      string            `yaml:"version"`
	Url          string            `yaml:"url"`
	DateFetched  string            `yaml:"date_fetched"`
	Count        int               `yaml:"trusted_certificates_count"`
	Certificates []observatoryCert `yaml:"trusted_certificates"`
}

type observatoryCert struct {
	SubjectName string `yaml:"subject_name"`
	Fingerprint string `yaml:"fingerprint"`
}

func writeObservatoryReport(meta Meta, certs []*x509.Certificate, cfg *Config) error {
	dateFetched := time.Now().Format(observatoryTimeFormat)
	count := len(certs)
	report := observatoryReport{
		Platform:    meta.Name,
		Version:     meta.Version,
		DateFetched: dateFetched,
		Count:       count,
	}
	obsCerts := make([]observatoryCert, count)
	for i := range obsCerts {
		cert := certs[i]
		obsCerts[i] = observatoryCert{
			SubjectName: certutil.StringifyPKIXName(cert.Subject),
			Fingerprint: certutil.GetHexSHA256Fingerprint(*cert),
		}
	}
	report.Certificates = obsCerts

	// write the report somewhere
	bs, err := yaml.Marshal(&report)
	if err != nil {
		return err
	}
	if cfg.Outfile != "" {
		return ioutil.WriteFile(cfg.Outfile, bs, 0644)
	}
	// write yaml to stdout
	_, err = os.Stdout.Write(bs)
	if err != nil {
		return err
	}
	return nil
}