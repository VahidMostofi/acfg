package endpointsagg

import (
	"testing"
)

func TestEndpointsAggregator_GetEndpointsResponseTimes(t *testing.T) {
	//viper.Set(constants.EndpointsAggregatorArgsURL, os.Getenv("INFLUXDB_URL"))
	//viper.Set(constants.EndpointsAggregatorArgsToken, os.Getenv("INFLUXDB_TOKEN"))
	//viper.Set(constants.EndpointsAggregatorArgsOrganization, os.Getenv("INFLUXDB_ORG"))
	//viper.Set(constants.EndpointsAggregatorArgsBucket, os.Getenv("INFLUXDB_BUCKET"))
	//
	//resourceFilters := map[string]map[string]interface{}{
	//	"login": {"URI_REGEX":"login*", "HTTP_METHOD":"POST"},
	//	"get-book": {"URI_REGEX":"books*", "HTTP_METHOD":"GET"},
	//	"edit-book": {"URI_REGEX":"books*", "HTTP_METHOD":"PUT"},
	//}
	//viper.Set(constants.EndpointsFilters, resourceFilters)
	//ea,err := NewEndpointsAggregator("influxdb")
	//if err != nil{
	//	t.Log(err)
	//	t.Fail()
	//	return
	//}
	//
	//type args struct {
	//	startTime  int64
	//	finishTime int64
	//}
	//tests := []struct {
	//	name    string
	//	args    args
	//}{
	//	{"does it work?", args{time.Now().Add(-3 * time.Minute).Unix(), time.Now().Add(-1 * time.Minute).Unix()}},
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//
	//		got, err := ea.GetEndpointsResponseTimes(tt.args.startTime, tt.args.finishTime)
	//		if err != nil{
	//			t.Errorf("GetEndpointsResponseTimes() error = %v",err)
	//		}
	//
	//		if len(got) != 3{
	//			t.Log("len(got) should be 3 but it is", len(got))
	//		}
	//
	//		for key, value := range got{
	//			m, err := value.GetMean()
	//			if err != nil{
	//				t.Log(err)
	//				t.Fail()
	//				return
	//			}
	//			p90, err := value.GetPercentile(90)
	//			if err != nil{
	//				t.Log(err)
	//				t.Fail()
	//				return
	//			}
	//			fmt.Println(key, m, p90, value.GetCount())
	//		}
	//	})
	//}
}
