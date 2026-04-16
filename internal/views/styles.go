package views

import "strings"

// Button Tailwind class constants.
const (
	btnBase    = "px-3 py-1.5 rounded-[6px] border text-[13px] font-semibold cursor-pointer inline-flex items-center gap-1 transition-[background,opacity] duration-150 active:scale-[.98] disabled:opacity-40 disabled:cursor-not-allowed"
	BtnStart   = btnBase + " bg-[rgba(22,163,74,.08)] text-[#15803d] border-[rgba(22,163,74,.3)] hover:bg-[rgba(22,163,74,.15)]"
	BtnStop    = btnBase + " bg-[rgba(220,38,38,.07)] text-[#dc2626] border-[rgba(220,38,38,.25)] hover:bg-[rgba(220,38,38,.14)]"
	BtnRestart = btnBase + " bg-[rgba(217,119,6,.07)] text-[#b45309] border-[rgba(217,119,6,.25)] hover:bg-[rgba(217,119,6,.14)]"
	BtnBrowser = btnBase + " bg-[#004494] text-white border-[#004494] hover:bg-[#003070]"

	// Badge class constants.
	badgeBase    = "inline-flex items-center text-[11px] font-bold px-2 py-0.5 rounded-full whitespace-nowrap tracking-[.02em]"
	BadgeOk      = badgeBase + " bg-[rgba(22,163,74,.1)] text-[#15803d] border border-[rgba(22,163,74,.3)]"
	BadgeStopped = badgeBase + " bg-[rgba(107,114,128,.1)] text-[#6b7280] border border-[rgba(107,114,128,.25)]"
)

// SafeID converts an arbitrary string into a CSS-safe identifier.
func SafeID(s string) string {
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '-'
	}, s)
}
