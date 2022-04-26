/*
 * Dinosaur Service Fleet Manager
 *
 * Dinosaur Service Fleet Manager APIs that are used by internal services e.g fleetshard operators.
 *
 * API version: 1.4.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package private

// ManagedDinosaurDetailsMetadata struct for ManagedDinosaurDetailsMetadata
type ManagedDinosaurDetailsMetadata struct {
	Name        string                                    `json:"name,omitempty"`
	Namespace   string                                    `json:"namespace,omitempty"`
	Annotations ManagedDinosaurDetailsMetadataAnnotations `json:"annotations,omitempty"`
}