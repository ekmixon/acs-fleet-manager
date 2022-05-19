/*
 * Dinosaur Service Fleet Manager
 *
 * Dinosaur Service Fleet Manager APIs that are used by internal services e.g fleetshard operators.
 *
 * API version: 1.4.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package private

// ManagedDinosaurList A list of ManagedDinosaur
type ManagedDinosaurList struct {
	Kind  string            `json:"kind"`
	Items []ManagedDinosaur `json:"items"`
}