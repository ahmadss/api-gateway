package certloader

import (
	"crypto/tls"
	"fmt"
	"time"
)

type CertManager struct {
	certFile string
	keyFile  string
	cert     *tls.Certificate
	lastMod  time.Time
}

func New(certFile, keyFile string) *CertManager {
	return &CertManager{
		certFile: certFile,
		keyFile:  keyFile,
	}
}

func (cm *CertManager) GetCertificate() (*tls.Certificate, error) {
	if time.Since(cm.lastMod) > 5*time.Minute {
		if err := cm.reload(); err != nil {
			return nil, err
		}
	}
	return cm.cert, nil
}

func (cm *CertManager) reload() error {
	cert, err := tls.LoadX509KeyPair(cm.certFile, cm.keyFile)
	if err != nil {
		return fmt.Errorf("加载证书失败: %v", err)
	}
	cm.cert = &cert
	cm.lastMod = time.Now()
	return nil
}
