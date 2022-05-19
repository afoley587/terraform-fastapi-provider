// resource_server.go
package main

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"log"
	"net/http"
)

func desDecrypt() *schema.Resource {
	return &schema.Resource{
		Create: resourceDesDecryptCreate,
		Read:   resourceDesDecryptRead,
		Update: resourceDesDecryptUpdate,
		Delete: resourceDesDecryptDelete,

		Schema: map[string]*schema.Schema{
			/* Our input to our resource */
			"ciphertext": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			/* Our computed attribute of our resource */
			"plaintext": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

type DesDecryptRequestBody struct {
	Ciphertext string `json:"ciphertext"`
}

type DesDecryptResponse struct {
	Id        string `json:"id"`
	Plaintext string `json:"plaintext"`
}

func resourceDesDecryptCreate(d *schema.ResourceData, m interface{}) error {
	/* Create the GoLang Object with the ciphertext from our resource */
	desRequestBody := DesDecryptRequestBody{
		Ciphertext: d.Get("ciphertext").(string),
	}

	/* Convert DesDecryptRequestBody to byte using Json.Marshal
	 * Ignoring error.
	 */
	body, _ := json.Marshal(desRequestBody)

	/* Actually send the payload!
	 * Make sure to send the JSON with the proper Content-Type Header
	 * Set the body to our JSON bytes
	 */
	resp, err := http.Post("http://127.0.0.1:8000/decrypt/des", "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Printf("[INFO] Could not initialize 'http.NewRequest' to 'http://localhost:8000/'.")
		return err
	}

	/* Reading the response body. This will be the JSON
	 * representation of the DESDecryptedResponse that we defined
	 * in the python/schemas/encryption.py file!
	 */
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("[INFO] Could not read response body.")
		return err
	}

	/* Deserialize the JSON/Byte response to our expected response */
	var respJsonDecoded DesDecryptResponse
	err = json.Unmarshal([]byte(string(respBody)), &respJsonDecoded)

	if err != nil {
		log.Printf("[INFO] Could not unmarshal JSON response.")
		return err
	}

	/* Save our outputs, id and plaintext
	 */
	d.SetId(respJsonDecoded.Id)
	d.Set("plaintext", respJsonDecoded.Plaintext)

	/* Best practice is just to read it from the API.
	 * of course, this does NOTHING in our case, but
	 * acts as a good sanity check
	 */
	resourceDesDecryptRead(d, m)
	return nil
}

func resourceDesDecryptRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDesDecryptUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceDesDecryptRead(d, m)
}

func resourceDesDecryptDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}
