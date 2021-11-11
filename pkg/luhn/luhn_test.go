package luhn

import (
	"testing"
)

func TestChecksum(t *testing.T) {
	tests := []struct {
		name    string
		number  string
		want    int
		wantErr bool
	}{
		{
			name:    "#1 Non valid number",
			number:  "Вот такой у нас номер",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "#2 Empty number",
			number:  "",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "#3 One-digit number",
			number:  "8",
			want:    3,
			wantErr: false,
		},
		{
			name:    "#4 Odd-digit number",
			number:  "123456789",
			want:    10 - (9+8+5+6+1+4+6+2+2)%10,
			wantErr: false,
		},
		{
			name:    "#5 Even-digit number",
			number:  "12345678",
			want:    10 - (7+7+3+5+8+3+4+1)%10,
			wantErr: false,
		},
		{
			name:    "#6 Number with 0 cheksum digit",
			number:  "778900746064935",
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Checksum(tt.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("Checksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Checksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestValidateCalculate performs first generate the number, that passes Luhn's check.
// Then that number is checked by Validate function.
func TestValidateCalculate(t *testing.T) {
	tests := []struct {
		name         string
		number       string
		wantErrCalc  bool
		wantValidate bool
	}{
		{
			name:         "#1 Wrong number",
			number:       "1234876O46",
			wantErrCalc:  true,
			wantValidate: false,
		},
		{
			name:         "#2 Empty number",
			number:       "",
			wantErrCalc:  true,
			wantValidate: false,
		},
		{
			name:         "#3 One-digit number",
			number:       "9",
			wantErrCalc:  false,
			wantValidate: true,
		},
		{
			name:         "#4 Odd-digit number",
			number:       "193247509384751938405193847019384751938476",
			wantErrCalc:  false,
			wantValidate: true,
		},
		{
			name:         "#5 Even-digit number",
			number:       "1932475109384751938405193847019384751938476",
			wantErrCalc:  false,
			wantValidate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generated, err := Calculate(tt.number)
			if (err == nil) == tt.wantErrCalc {
				t.Errorf("Calculate() error = %v, wantErr = %v", err, tt.wantErrCalc)
			}
			if got := Validate(generated); got != tt.wantValidate {
				t.Errorf("Validate() = %v, want %v", got, tt.wantValidate)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	validnumber := "7789007064617936" // This is my credit card number. You may buy something as a gift from me:))
	invalidnumber := "123456789"

	// True case
	if !Validate(validnumber) {
		t.Errorf("%s must pass the check.", validnumber)
	}

	// False case
	if Validate(invalidnumber) {
		t.Errorf("%s must not pass the check.", invalidnumber)
	}
}
