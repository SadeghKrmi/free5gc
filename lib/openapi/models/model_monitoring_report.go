/*
 * Nudm_EE
 *
 * Nudm Event Exposure Service
 *
 * API version: 1.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

import (
	"time"
)

type MonitoringReport struct {
	ReferenceId int32      `json:"referenceId" yaml:"referenceId" bson:"referenceId" mapstructure:"ReferenceId"`
	EventType   EventType  `json:"eventType" yaml:"eventType" bson:"eventType" mapstructure:"EventType"`
	Report      *Report    `json:"report,omitempty" yaml:"report" bson:"report" mapstructure:"Report"`
	Gpsi        string     `json:"gpsi,omitempty" yaml:"gpsi" bson:"gpsi" mapstructure:"Gpsi"`
	TimeStamp   *time.Time `json:"timeStamp" yaml:"timeStamp" bson:"timeStamp" mapstructure:"TimeStamp"`
}
