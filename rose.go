package rose

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

func formatJSON(reader io.Reader, writer io.Writer) error {
	var decoded map[string]interface{}
	err := json.NewDecoder(reader).Decode(&decoded)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "\t")
	return encoder.Encode(decoded)
}

func formatYAML(reader io.Reader, writer io.Writer) error {
	var decoded map[string]interface{}
	err := xml.NewDecoder(reader).Decode(&decoded)
	if err != nil {
		return err
	}
	encoder := xml.NewEncoder(writer)
	return encoder.Encode(decoded)
}

//func goldenJSON(t *testing.T, golden string, reader io.Reader) {
//	if update {
//		file, err := os.OpenFile(golden, os.O_WRONLY, os.ModeExclusive)
//		require.NoError(t, err)
//		err = formatJSON(reader, file)
//		require.NoError(t, err)
//	} else {
//		got, err := ioutil.ReadAll(reader)
//		require.NoError(t, err)
//		expected, err := ioutil.ReadFile(golden)
//		require.NoError(t, err)
//		require.JSONEq(t, string(expected), string(got), "request does not match expected")
//	}
//}
