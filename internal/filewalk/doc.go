// Package filewalk owns repository file discovery for style checks.
//
// It applies repository scopes, global exclusions, generated-file detection,
// binary probing, deterministic sorting, and path normalisation before callers
// apply rule-specific policy.
package filewalk
