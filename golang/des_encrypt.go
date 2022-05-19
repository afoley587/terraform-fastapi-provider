// des_encrypt.go
package main

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"log"
	"net/http"
)

func desEncrypt() *schema.Resource {
	return &schema.Resource{
		Create: resourceDesEncryptCreate,
		Read:   resourceDesEncryptRead,
		Update: resourceDesEncryptUpdate,
		Delete: resourceDesEncryptDelete,

		Schema: map[string]*schema.Schema{
			/* Our input to our resource */
			"plaintext": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			/* Our computed attribute of our resource */
			"ciphertext": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

type DesEncryptRequestBody struct {
	Plaintext string `json:"plaintext"`
}

type DesEncryptResponse struct {
	Id         string `json:"id"`
	Ciphertext string `json:"ciphertext"`
}

func resourceDesEncryptCreate(d *schema.ResourceData, m interface{}) error {
	/* Create the GoLang Object with the ciphertext from our resource */
	log.Printf("[INFO] Starting call to http://127.0.0.1:8000/")
	desRequestBody := DesEncryptRequestBody{
		Plaintext: d.Get("plaintext").(string),
	}

	/* Convert DesEncryptRequestBody to byte using Json.Marshal
	 * Ignoring error.
	 */
	body, _ := json.Marshal(desRequestBody)

	/* Actually send the payload!
	 * Make sure to send the JSON with the proper Content-Type Header
	 * Set the body to our JSON bytes
	 */
	resp, err := http.Post("http://127.0.0.1:8000/encrypt/des", "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Printf("[INFO] Could not initialize 'http.NewRequest' to 'http://localhost:8000/'.")
		return err
	}

	/* Reading the response body. This will be the JSON
	 * representation of the DESEncryptedResponse that we defined
	 * in the python/schemas/encryption.py file!
	 */
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("[INFO] Could not read response body.")
		return err
	}

	/* Deserialize the JSON/Byte response to our expected response */
	var respJsonDecoded DesEncryptResponse
	err = json.Unmarshal([]byte(string(respBody)), &respJsonDecoded)

	if err != nil {
		log.Printf("[INFO] Could not unmarshal JSON response.")
		return err
	}

	/* Save our outputs, id and plaintext
	 */
	d.SetId(respJsonDecoded.Id)
	d.Set("ciphertext", respJsonDecoded.Ciphertext)
	/* Best practice is just to read it from the API.
	 * of course, this does NOTHING in our case, but
	 * acts as a good sanity check
	 */
	resourceDesEncryptRead(d, m)
	return nil
}

func resourceDesEncryptRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDesEncryptUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceDesEncryptRead(d, m)
}

func resourceDesEncryptDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}
