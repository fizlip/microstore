package models 

var MockUserAnalytics = []UserAnalytics{
	{
		CustomerID:   "filip",
		Timestamp:    "2024-09-24",
		PaidRequests: 123,
		FreeRequests: 502,
		Products: map[string]ProductData{
			"p_001": {ProductID: "123", ProductName: "name1", Price: 0.011, NumRequests: 18},
			"p_002": {ProductID: "124", ProductName: "name2", Price: 0.0011, NumRequests: 21},
			"p_003": {ProductID: "125", ProductName: "name3", Price: 0.021, NumRequests: 14},
			"p_004": {ProductID: "126", ProductName: "name3", Price: 0.01, NumRequests: 110},
			"p_005": {ProductID: "127", ProductName: "name4", Price: 0.1, NumRequests: 110},
			"p_006": {ProductID: "128", ProductName: "name6", Price: 0.0, NumRequests: 101},
		},
	},
	{
		CustomerID:   "filip",
		Timestamp:    "2024-09-23",
		PaidRequests: 328,
		FreeRequests: 1021,
		Products: map[string]ProductData{
			"p_001": {ProductID: "123", ProductName: "name1", Price: 0.011, NumRequests: 81},
			"p_002": {ProductID: "124", ProductName: "name2", Price: 0.0011, NumRequests: 210},
			"p_003": {ProductID: "125", ProductName: "name3", Price: 0.021, NumRequests: 121},
			"p_004": {ProductID: "126", ProductName: "name3", Price: 0.01, NumRequests: 434},
			"p_005": {ProductID: "127", ProductName: "name4", Price: 0.1, NumRequests: 842},
			"p_006": {ProductID: "128", ProductName: "name6", Price: 0.0, NumRequests: 921},
		},
	},
	{
		CustomerID:   "filip",
		Timestamp:    "2024-09-23",
		PaidRequests: 223,
		FreeRequests: 1092,
		Products: map[string]ProductData{
			"p_001": {ProductID: "123", ProductName: "name1", Price: 0.011, NumRequests: 18},
			"p_002": {ProductID: "124", ProductName: "name2", Price: 0.0011, NumRequests: 21},
			"p_003": {ProductID: "125", ProductName: "name3", Price: 0.021, NumRequests: 14},
			"p_004": {ProductID: "126", ProductName: "name3", Price: 0.01, NumRequests: 110},
			"p_005": {ProductID: "127", ProductName: "name4", Price: 0.1, NumRequests: 110},
			"p_006": {ProductID: "128", ProductName: "name6", Price: 0.0, NumRequests: 101},
		},
	},
	{
		CustomerID:   "filip",
		Timestamp:    "2024-09-22",
		PaidRequests: 291,
		FreeRequests: 891,
		Products: map[string]ProductData{
			"p_001": {ProductID: "123", ProductName: "name1", Price: 0.011, NumRequests: 10},
			"p_002": {ProductID: "124", ProductName: "name2", Price: 0.0011, NumRequests: 19},
			"p_003": {ProductID: "125", ProductName: "name3", Price: 0.021, NumRequests: 14},
			"p_004": {ProductID: "126", ProductName: "name3", Price: 0.01, NumRequests: 192},
			"p_005": {ProductID: "127", ProductName: "name4", Price: 0.1, NumRequests: 121},
			"p_006": {ProductID: "128", ProductName: "name6", Price: 0.0, NumRequests: 111},
		},
	},
}
