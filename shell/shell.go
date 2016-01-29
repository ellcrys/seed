// A shell represents a container of information
// that is considered valueable.
package shell

import (
	"encoding/json"
	"errors"
	"fmt"
	// "reflect"
)

type Shell struct {
	Signatures map[string]interface{} 		`json:"signatures"`
	Meta map[string]interface{}				`json:"meta"`
	Ownership map[string]interface{} 		`json:"ownership"`
	Embeds []map[string]interface{} 		`json:"embeds"`
	Attributes map[string]interface{}		`json:"attributes"`
}

// Initialize a shell
func initialize(shell *Shell) *Shell {
	shell.Signatures = make(map[string]interface{})
	shell.Meta = make(map[string]interface{})
	shell.Ownership = make(map[string]interface{})
	shell.Embeds = []map[string]interface{}{}
	shell.Attributes = make(map[string]interface{})
	return shell
}

// Create a shell.The new shell is immediately signed using the issuer's private key
func Create(meta map[string]interface{}, issuerPrivateKey string) (*Shell, error) {

	shell := initialize(&Shell{})

	// validate meta
	if err := ValidateMetaBlock(meta); err != nil {
    	return &Shell{}, err
    }

    // set shell Meta field and create a meta signature
	shell.Meta = meta
	metaSignature, err := shell.Sign("meta", issuerPrivateKey)
	if err != nil {
		return &Shell{}, err
	}

	// added meta signature in signatures map
	shell.Signatures["meta"] = metaSignature
	return shell, nil
}

// Creates a shell from a map
func loadMap(data map[string]interface{}) (*Shell, error) {

	var shell = &Shell{}

	// add signatures
    if signatures := data["signatures"]; signatures != nil {
    	switch val := signatures.(type) {
    	case map[string]interface{}:
    		shell.Signatures = val
    		break;
    	default:
    		return &Shell{}, errors.New("`signatures` block has invalid value type. Expects JSON object")
    	}
    }

    // add meta
    if meta := data["meta"]; meta != nil {
    	switch val := meta.(type) {
    	case map[string]interface{}:
    		shell.Meta = val
    		break;
    	default:
    		return &Shell{}, errors.New("`meta` block has invalid value type. Expects JSON object")
    	}
    }

    // add ownership
    if ownership := data["ownership"]; ownership != nil {
    	switch val := ownership.(type) {
    	case map[string]interface{}:
    		shell.Ownership = val
    		break;
    	default:
    		return &Shell{}, errors.New("`ownership` block has invalid value type. Expects JSON object")
    	}
    }

    // // add attributes
    if attributes := data["attributes"]; attributes != nil {
    	switch val := attributes.(type) {
    	case map[string]interface{}:
    		shell.Attributes = val
    		break;
    	default:
    		return &Shell{}, errors.New("`attributes` block has invalid value type. Expects JSON object")
    	}
    }

    return shell, nil
}


// Create a shell from a json string by converting
// it to a map and then used to load a new shell instance
func LoadJSON(jsonStr string) (*Shell, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
        return &Shell{}, errors.New("unable to parse json string");
    }
    if err := Validate(data); err != nil {
    	return &Shell{}, err
    }
	return loadMap(data)  
}

// Sign any shell block by creating a canonical string representation
// of the block value and signing with the issuer's private key
func(self *Shell) Sign(blockName string, privateKey string) (string, error) {
	switch blockName {
	case "meta":
		canonicalMap := GetCanonicalMapString(self.Meta)
		signer, err := ParsePrivateKey([]byte(privateKey))
		if err != nil {
			return "", errors.New(fmt.Sprintf("Private Key Error: %v", err))
		}
		signature, err := signer.Sign([]byte(canonicalMap))
		if err != nil {
			return "", errors.New(fmt.Sprintf("Signature Error: %v", err))
		}
		return signature, nil
	default:
		return "", errors.New("block unknown")
	}
}

// return shell as raw JSON string
func(self *Shell) JSON() string {
	bs, _ := json.Marshal(&self)
	return string(bs)
}

// return shell as a base64 json encoded string
func(self *Shell) Encode() string {
	jsonStr := self.JSON()
	return ToBase64([]byte(jsonStr))
}

