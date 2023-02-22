// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package models contain request and response structures used for the SAPControl Web Service Interface
package models // import "github.com/observiq/observiq-otel-collector/receiver/sapnetweaverreceiver/internal/models"

import "encoding/xml"

// StateColor holds the string state color value
type StateColor string

// StateColorCode holds the numerical state color value
type StateColorCode int

// StateColor constant values
const (
	StateColorGray       StateColor     = "SAPControl-GRAY"
	StateColorGreen      StateColor     = "SAPControl-GREEN"
	StateColorYellow     StateColor     = "SAPControl-YELLOW"
	StateColorRed        StateColor     = "SAPControl-RED"
	StateColorCodeGray   StateColorCode = 1
	StateColorCodeGreen  StateColorCode = 2
	StateColorCodeYellow StateColorCode = 3
	StateColorCodeRed    StateColorCode = 4
)

// GetAlertTreeResponse is an xml response struct
type GetAlertTreeResponse struct {
	XMLName   xml.Name     `xml:"urn:SAPControl GetAlertTreeResponse"`
	AlertNode []*AlertNode `xml:"tree>item,omitempty" json:"tree>item,omitempty"`
}

// AlertNode is an xml response struct
type AlertNode struct {
	Name        string     `xml:"name,omitempty" json:"name,omitempty"`
	ActualValue StateColor `xml:"ActualValue,omitempty" json:"ActualValue,omitempty"`
	Description string     `xml:"description,omitempty" json:"description,omitempty"`
}

// GetInstancePropertiesResponse is an xml response struct
type GetInstancePropertiesResponse struct {
	XMLName    xml.Name            `xml:"urn:SAPControl GetInstancePropertiesResponse"`
	Properties []*InstanceProperty `xml:"properties>item,omitempty" json:"properties>item,omitempty"`
}

// InstanceProperty is an xml response struct
type InstanceProperty struct {
	Property string `xml:"property,omitempty" json:"property,omitempty"`
	Value    string `xml:"value,omitempty" json:"value,omitempty"`
}

// CurrentSapInstance is a struct to hold the instance xml instant property response values
type CurrentSapInstance struct {
	SID      string
	Number   int32
	Name     string
	Hostname string
}

// GetSystemInstanceListResponse is an xml response struct
type GetSystemInstanceListResponse struct {
	XMLName  xml.Name            `xml:"urn:SAPControl GetSystemInstanceListResponse"`
	Instance *ArrayOfSAPInstance `xml:"instance,omitempty" json:"instance,omitempty"`
}

// ArrayOfSAPInstance is an xml response struct
type ArrayOfSAPInstance struct {
	Item []*SAPInstance `xml:"item,omitempty" json:"item,omitempty"`
}

// SAPInstance is an xml response struct
type SAPInstance struct {
	Hostname      string     `xml:"hostname,omitempty" json:"hostname,omitempty"`
	InstanceNr    int32      `xml:"instanceNr,omitempty" json:"instanceNr,omitempty"`
	StartPriority string     `xml:"startPriority,omitempty" json:"startPriority,omitempty"`
	Features      string     `xml:"features,omitempty" json:"features,omitempty"`
	Dispstatus    StateColor `xml:"dispstatus,omitempty" json:"dispstatus,omitempty"`
}

// GetProcessListResponse is an xml response struct
type GetProcessListResponse struct {
	XMLName xml.Name          `xml:"urn:SAPControl GetProcessListResponse"`
	Process *ArrayOfOSProcess `xml:"process,omitempty" json:"process,omitempty"`
}

// ArrayOfOSProcess is an xml response struct
type ArrayOfOSProcess struct {
	Item []*OSProcess `xml:"item,omitempty" json:"item,omitempty"`
}

// OSProcess is an xml response struct
type OSProcess struct {
	Name        string      `xml:"name,omitempty" json:"name,omitempty"`
	Description string      `xml:"description,omitempty" json:"description,omitempty"`
	Dispstatus  *StateColor `xml:"dispstatus,omitempty" json:"dispstatus,omitempty"`
	Textstatus  string      `xml:"textstatus,omitempty" json:"textstatus,omitempty"`
	Starttime   string      `xml:"starttime,omitempty" json:"starttime,omitempty"`
	Elapsedtime string      `xml:"elapsedtime,omitempty" json:"elapsedtime,omitempty"`
	Pid         int32       `xml:"pid,omitempty" json:"pid,omitempty"`
}

// GetQueueStatisticResponse is an xml response struct
type GetQueueStatisticResponse struct {
	XMLName xml.Name                 `xml:"urn:SAPControl GetQueueStatisticResponse"`
	Queue   *ArrayOfTaskHandlerQueue `xml:"queue,omitempty" json:"queue,omitempty"`
}

// ArrayOfTaskHandlerQueue is an xml response struct
type ArrayOfTaskHandlerQueue struct {
	Item []*TaskHandlerQueue `xml:"item,omitempty" json:"item,omitempty"`
}

// TaskHandlerQueue is an xml response struct
type TaskHandlerQueue struct {
	Typ  string `xml:"Typ,omitempty" json:"Typ,omitempty"`
	Now  int32  `xml:"Now,omitempty" json:"Now,omitempty"`
	High int32  `xml:"High,omitempty" json:"High,omitempty"`
	Max  int32  `xml:"Max,omitempty" json:"Max,omitempty"`
}

// ABAPGetSystemWPTableResponse is an xml response struct
type ABAPGetSystemWPTableResponse struct {
	XMLName     xml.Name                  `xml:"urn:SAPControl ABAPGetSystemWPTableResponse"`
	Workprocess *ArrayOfSystemWorkProcess `xml:"workprocess,omitempty" json:"workprocess,omitempty"`
}

// ArrayOfSystemWorkProcess is an xml response struct
type ArrayOfSystemWorkProcess struct {
	Item []*SystemWorkProcess `xml:"item,omitempty" json:"item,omitempty"`
}

// SystemWorkProcess is an xml response struct
type SystemWorkProcess struct {
	Instance string `xml:"Instance,omitempty" json:"Instance,omitempty"`
	No       int32  `xml:"No,omitempty" json:"No,omitempty"`
	Typ      string `xml:"Typ,omitempty" json:"Typ,omitempty"`
	Pid      int32  `xml:"Pid,omitempty" json:"Pid,omitempty"`
	Status   string `xml:"Status,omitempty" json:"Status,omitempty"`
	Reason   string `xml:"Reason,omitempty" json:"Reason,omitempty"`
	Start    string `xml:"Start,omitempty" json:"Start,omitempty"`
	Err      string `xml:"Err,omitempty" json:"Err,omitempty"`
	Sem      string `xml:"Sem,omitempty" json:"Sem,omitempty"`
	CPU      string `xml:"Cpu,omitempty" json:"Cpu,omitempty"`
	Time     string `xml:"Time,omitempty" json:"Time,omitempty"`
	Program  string `xml:"Program,omitempty" json:"Program,omitempty"`
	Client   string `xml:"Client,omitempty" json:"Client,omitempty"`
	User     string `xml:"User,omitempty" json:"User,omitempty"`
	Action   string `xml:"Action,omitempty" json:"Action,omitempty"`
	Table    string `xml:"Table,omitempty" json:"Table,omitempty"`
}

// GetRequestLogFileResponse is an xml response struct
type GetRequestLogFileResponse struct {
	XMLName xml.Name  `xml:"urn:SAPControl GetRequestLogFileResponse"`
	Content []*string `xml:"content,omitempty" json:"content,omitempty"`
}

// OSExecuteResponse is an xml response struct
type OSExecuteResponse struct {
	XMLName  xml.Name       `xml:"urn:SAPControl OSExecuteResponse"`
	Exitcode int32          `xml:"exitcode,omitempty" json:"exitcode,omitempty"`
	Pid      int32          `xml:"pid,omitempty" json:"pid,omitempty"`
	Lines    *ArrayOfString `xml:"lines,omitempty" json:"lines,omitempty"`
}

// ArrayOfString is an xml response struct
type ArrayOfString struct {
	Item []string `xml:"item,omitempty" json:"item,omitempty"`
}

// EnqGetStatisticResponse is an xml response struct
type EnqGetStatisticResponse struct {
	XMLName       xml.Name `xml:"urn:SAPControl EnqStatistic"`
	LocksNow      *int32   `xml:"locks-now,omitempty" json:"locks-now,omitempty"`
	LocksHigh     *int32   `xml:"locks-high,omitempty" json:"locks-high,omitempty"`
	LocksMax      *int32   `xml:"locks-max,omitempty" json:"locks-max,omitempty"`
	EnqueueErrors *int64   `xml:"enqueue-errors,omitempty" json:"enqueue-errors,omitempty"`
	DequeueErrors *int64   `xml:"dequeue-errors,omitempty" json:"dequeue-errors,omitempty"`
	LockTime      *float64 `xml:"lock-time,omitempty" json:"lock-time,omitempty"`
	LockWaitTime  *float64 `xml:"lock-wait-time,omitempty" json:"lock-wait-time,omitempty"`
}
