// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import "math"

func ComputeMean(values []float64) float64 {
	var sum float64

	for _, value := range values {
		sum = sum + value
	}

	return (sum / float64(len(values)))
}

func ComputeStdDeviation(values []float64) float64 {
	var (
		mean         float64
		vp           []float64
		stdDiv, temp float64
	)

	vp = make([]float64, len(values))

	/*Calculate mean of the data points*/
	mean = ComputeMean(values)
	/*Calculate standard deviation of individual data points*/
	for i, v := range values {
		temp = v - mean
		vp[i] = (temp * temp)
	}

	/* Finally, Compute standard Deviation of data points
	 * by taking mean of individual std. Deviation.
	 */
	stdDiv = ComputeMean(vp)
	return float64(math.Sqrt(float64(stdDiv)))
}

