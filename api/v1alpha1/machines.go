package v1alpha1

import (
	"fmt"
	"net"

	"github.com/go-playground/validator/v10"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MAC is a MAC address.
type MAC net.HardwareAddr

// UnmarshalYAML unmarshals a MAC address from a string.
func (m *MAC) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var mac string

	if err := unmarshal(&mac); err != nil {
		return fmt.Errorf("failed to unmarshal MAC address: %w", err)
	}

	hw, err := net.ParseMAC(mac)
	if err != nil {
		return fmt.Errorf("failed to parse MAC address: %w", err)
	}

	*m = MAC(hw)

	return nil
}

// MarshalYAML marshals a MAC address to a string.
func (m MAC) MarshalYAML() (interface{}, error) {
	return net.HardwareAddr(m).String(), nil
}

// Interface describes a network interface of a Machine.
type Interface struct {
	// MAC is the MAC address of the interface.
	MAC MAC `json:"mac" validate:"required"`
}

// MachineSpecHardware defines the hardware configuration of a Machine.
type MachineSpecHardware struct {
	// Vendor is the manufacturer of the machine.
	Vendor string `json:"vendor" validate:"required"`
	// Model is the model of the machine.
	Model  string `json:"model" validate:"required"`
}

// MachineSpec defines the desired state of Machine.
type MachineSpec struct {
	// Hardware is the hardware configuration of the machine.
	Hardware MachineSpecHardware `json:"hardware" validate:"required"`
	// Interfaces describes the network interfaces of the machine.
	Interfaces []Interface `json:"interfaces" validate:"required,dive"`
}

// Machine represents a physical machine in the homelab.
type Machine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the machine.
	Spec   MachineSpec   `json:"spec" validate:"required"`
}

// Validate validates the Machine.
func (m *Machine) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(m); err != nil {
		return fmt.Errorf("failed to validate Machine: %w", err)
	}

	return nil
}

// MachineList contains a list of Machine.
type MachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of Machine objects in the list.
	Items           []Machine `json:"items"`
}
