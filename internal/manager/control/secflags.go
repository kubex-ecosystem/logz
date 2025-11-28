package control

import (
	"fmt"
	"sort"
	"strings"
)

/* ========= SECURITY FLAGS ========= */

// SecFlag representa opções de segurança aplicáveis a uma rota, handler, etc.
type SecFlag uint32

const (
	SecNone         SecFlag = 0
	SecAuth         SecFlag = 1 << iota // autenticação/autorização
	SecSanitize                         // sanitize em params/headers
	SecSanitizeBody                     // sanitize no body
)

// Has retorna true se TODOS os bits de mask estiverem presentes.
func (f SecFlag) Has(mask SecFlag) bool {
	return f&mask == mask
}

// With adiciona um ou mais flags.
func (f SecFlag) With(mask SecFlag) SecFlag {
	return f | mask
}

// Without remove um ou mais flags.
func (f SecFlag) Without(mask SecFlag) SecFlag {
	return f &^ mask
}

// ordem determinística para log/telemetria
var secOrder = []struct {
	name string
	flag SecFlag
}{
	{"auth", SecAuth},
	{"sanitize", SecSanitize},
	{"sanitize_body", SecSanitizeBody},
}

// String retorna uma representação textual dos flags ativos.
// Ex: "auth|sanitize_body" ou "none" se não houver nenhum setado.
func (f SecFlag) String() string {
	if f == SecNone {
		return "none"
	}
	var parts []string
	for _, it := range secOrder {
		if f.Has(it.flag) {
			parts = append(parts, it.name)
		}
	}
	if len(parts) == 0 {
		return fmt.Sprintf("unknown(0x%X)", uint32(f))
	}
	sort.Strings(parts)
	return strings.Join(parts, "|")
}

/* ======== MAP LEGADO -> FLAGS ======== */

// FromLegacyMap converte mapa legado (ex: de config JSON) para flags.
// Mantém compat com chaves antigas: "validateAndSanitize", "validateAndSanitizeBody".
func FromLegacyMap(m map[string]bool) SecFlag {
	if m == nil {
		return SecNone
	}
	var f SecFlag
	if m["auth"] {
		f |= SecAuth
	}
	if m["sanitize"] {
		f |= SecSanitize
	}
	if m["sanitize_body"] || m["validateAndSanitizeBody"] {
		f |= SecSanitizeBody
	}
	// compat antiga: se só tinha validateAndSanitize, assume sanitize geral.
	if m["validateAndSanitize"] {
		f |= SecSanitize
	}
	return f
}
