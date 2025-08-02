package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
	"strings"
	"testing"
)

func TestToArgs(t *testing.T) {
	//
	t.Run("string value", func(t *testing.T) {
		result := ToArgs("test")
		expected := &structpb.Value{
			Kind: &structpb.Value_StringValue{
				StringValue: "test",
			},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("bool value", func(t *testing.T) {
		result := ToArgs(true)
		expected := &structpb.Value{
			Kind: &structpb.Value_BoolValue{
				BoolValue: true,
			},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("float value", func(t *testing.T) {
		result := ToArgs(3.14)
		expected := &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: 3.14,
			},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("slice value", func(t *testing.T) {
		input := []interface{}{"test", true, 3.14}
		result := ToArgs(input)

		expected := &structpb.Value{
			Kind: &structpb.Value_ListValue{
				ListValue: &structpb.ListValue{
					Values: []*structpb.Value{
						{Kind: &structpb.Value_StringValue{StringValue: "test"}},
						{Kind: &structpb.Value_BoolValue{BoolValue: true}},
						{Kind: &structpb.Value_NumberValue{NumberValue: 3.14}},
					},
				},
			},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("nested slice value", func(t *testing.T) {
		input := []interface{}{
			[]interface{}{"nested", true},
			"test",
		}
		result := ToArgs(input)

		expected := &structpb.Value{
			Kind: &structpb.Value_ListValue{
				ListValue: &structpb.ListValue{
					Values: []*structpb.Value{
						{
							Kind: &structpb.Value_ListValue{
								ListValue: &structpb.ListValue{
									Values: []*structpb.Value{
										{Kind: &structpb.Value_StringValue{StringValue: "nested"}},
										{Kind: &structpb.Value_BoolValue{BoolValue: true}},
									},
								},
							},
						},
						{Kind: &structpb.Value_StringValue{StringValue: "test"}},
					},
				},
			},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("nil value", func(t *testing.T) {
		var input interface{} = nil
		result := ToArgs(input)
		assert.Nil(t, result)
	})
}

func TestCensorString(t *testing.T) {
	testCases := []string{
		"This is a normal string.",
		"My API key is sk-live-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		"Another key: pk-test-yyyyyyyyyyyyyyyyyyyyyyyyyyyyy",
		"The access token is Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		"client_secret: SomeReallyLongAndComplexSecretValueHere123ABC",
		"Generic token: 0123456789abcdef0123456789abcdef",
		"No secrets here, just some text.",
		"Multiple secrets: key1=sk-abc, key2=pk-def, token=Bearer xyz.123.abc",
		"Short key: abcdefg", // This might not be caught by some broad regexes
		"An ID that looks like a hex key but isn't: 1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d",
		"https://graph.facebook.com/v22.0/act_716524384010035/adsets?fields=id,name,status,effective_status,optimization_goal,destination_type,campaign_id,account_id,targeting{geo_locations{countries}}&access_token=EAAFZAo0S0PR0BOwjOnr821FhZAWpb5vkGAEpZCEItOiJXdNYs7dHghCerFBZCHGP5nzavk33VYhBlTBwr9vo2jG5fb3fMnnaQA88SvrHnzdfKRf4BcEAYhGQ8UfPc0tPXeyLwq6xUMcAZA834p3t0l0Alhf9zZCq79xISZB9a1Th1DksWV36ivjKOdc9r3D&limit=25&after=QVFIUmxCb0FXQlA4d2MxOTdXejJNa3pQWVVETC1oSWE5RUV3RGtrNzdwWE02S3Vib0VvaFFfVGp1TDRUNHdvZAEZAucVhqd3pZAS180QXVYeGRpTUhjcVRqQXh3",
	}

	for _, tc := range testCases {
		censored, found := DetectAndCensorAPIKeysTokens(tc)
		fmt.Printf("Original: \"%s\"\n", tc)
		fmt.Printf("Censored: \"%s\" (Found Sensitive: %t)\n", censored, found)
		fmt.Println(strings.Repeat("-", 80))
	}
}
