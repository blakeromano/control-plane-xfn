package s3

var publicRule = []map[string]interface{}{
	map[string]interface{}{
		"status": "Disabled",
		"expiration": []interface{}{
			map[string]interface{}{
				"days": 9999,
			},
		},
	},
}

var internalRule = []interface{}{
	map[string]interface{}{
		"status": "Enabled",
		"expiration": []interface{}{
			map[string]interface{}{
				"days": 1095,
			},
		},
	},
}

var confidentialRule = []interface{}{
	map[string]interface{}{
		"status": "Enabled",
		"expiration": []interface{}{
			map[string]interface{}{
				"days": 730,
			},
		},
	},
}