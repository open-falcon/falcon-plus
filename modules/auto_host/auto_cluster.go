
func genAggregator(){
	endpointCounterList, err := getEndpointCounters()
	for _,endpointCounter := range endpointCounterList{
		genAggr(endpointCounter)
	}

}

func getEndpointCounters(){

}
func genAggr(endpointCounter) {
	grpId, grpEndpointName := getGrpinfo(endpointCounter.EndpointID)
	numberator := getNumberator(endpointCounter.Counter)
	denominator:= getDenominator(endpointCounter.Type)
	metric := getMetric(endpointCounter.Counter)
	tags := getTags(endpointCounter.Counter)
	dstype := getDstype(endpointCounter.Type)
	cluster  := Cluster  {
		GrpId  :grpId,
		Numerator  :numberator,
		Denominator :denominator,
		Endpoint  :grpEndpointName,
		Metric     :metric,
		Tags    :tags,
		DsType   :dstype,
		Step       :endpointCounter.Step,
		Creator    : "bot",
	}
	
}

func getEndpointName(endpointID int64){

}

func getGrpName(endpointID int64){

}
func getGrpId(endpointID int64){

}

func getGrpinfo(endpointID int64){
	name := getEndpointName(endpointID)
	grpNmae := getGrpName(name)
	grpId :=  getGrpId(grpNmae)
	return grpNmae,grpId
}


func getNumberator(endpointCounter.Counter){

}

func getDenominator(endpointCounter.Type){

}

func getMetric(endpointCounter.Counter){

}

func getTags(endpointCounter.Counter){

}

func getDstype(endpointCounter.Type){

}
