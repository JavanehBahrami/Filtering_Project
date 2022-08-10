/*
	This code implements a filtering serving named http_server 
	Author: Saeedeh (Javaneh) Bahrami
	Email: Bahramisaeede@gmail.com
*/
package main


import (
	"fmt"
	"errors"
	"os"
	"time"
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"io"
	"encoding/json"
	"log"
	"bufio"
	"github.com/mitchellh/mapstructure"
	"reflect"
   )


// File_name is the name of a text file inorder to save the outputs
const File_name = "responses_file.txt"

//assign a free port number
const PortNum = "8080"

//number of feilds in each input request is 4, including {x, y, width, height}
const Num_Feilds = 4


/*
	Predefined keys name in request data
	including {main{x, y, width, height}, input[{x, y, width, height}, ...]}
*/
var allowedExtKeys = []string{"main", "input"}
var allowedIntKeys = []string{"x", "y", "width", "height"}



type InputData struct {
	X int `json:"x"`
	Y int `json:"y"`
	Width int `json:"width"`
	Height int `json:"height"`
}


type MainData struct {
	X int `json:"x"`
	Y int `json:"y"`
	Width int `json:"width"`
	Height int `json:"height"`
}


// a struct to save the json objects request including{MainData, InputData}
type ReqData struct {
	Main MainData `json:"main"`
	Input []InputData `json:"input"`

}

// define a struct to save the json objects response
type SavedData struct {
	X int `json:"x"`
	Y int `json:"y"`
	Width int `json:"width"`
	Height int `json:"height"`
	Time string `json:"time"`
}


/* 
   This function handles GET requset
   Using this command to run it: curl -X GET -s localhost:8080 
*/
func handleGet(w http.ResponseWriter, req *http.Request) {
	file, err := os.Open(File_name)
    defer file.Close()

    if err != nil {
		log.Println("[INFO] Can not open File_name, %s", err)
		os.Exit(1)
    }else{
		reader := bufio.NewReader(file)
		io.WriteString(w, "[\n")
		var line string
		for {
			line, err = reader.ReadString('\n')
			io.WriteString(w, line)

			if err != nil {
				break
			}
		}

		if err != io.EOF {
			log.Println("[INFO] Failed reading File_name, %v\n", err)
		}
		io.WriteString(w, "]\n")
		file.Close()
	}

}


/* 
   This function filterd thoes rectangles which do not have any intersection 
   with the main coordinate box
*/
func check_intersection(main_data MainData, input_data InputData)string{
	x1 := main_data.X
	y1 := main_data.Y
	x2 := main_data.Width + x1
	y2 := main_data.Height + y1

	input_x1 := input_data.X
	input_y1 := input_data.Y
	input_x2 := input_data.Width + input_x1
	input_y2 := input_data.Height + input_y1

	if input_x2 > x1 && input_y2 > y1 && input_x2 < x2 && input_y2 < y2{
		return "valid"
	} else if input_x2 > x1 && input_y1 > y1 && input_x2 < x2 && input_y1 < y2{
		return "valid"
	} else if (input_x1 > x1 && input_y1 > y1 && input_x1 < x2 && input_y1 < y2){
		return "valid"
	} else if input_x1 > x1 && input_y2 > y1 && input_x1 < x2 && input_y2 < y2{
		return "valid"
	} else{
		return "filtered"
	}
}


/*
	This function processes the request data and
	saved them if they can pass the filter
*/
func parse_value(main_data MainData, input_data []InputData){
	Resp_file, err := os.OpenFile(File_name,
								os.O_APPEND|os.O_CREATE|os.O_WRONLY,
								0644)

	if err != nil {
		log.Println("[Info] Failed creating or openning Resp_file: %s", err)
		log.Println("[INFO] Try entering a correct file path and run the web server again!")
		os.Exit(1)
	}
	defer Resp_file.Close()
	datawriter := bufio.NewWriter(Resp_file)
	status_save := false

	var f_slice []string

	for i :=0; i < len(input_data); i++{
		// filtering_status:= check_intersection(main_data, input_data[i])
		input_x := input_data[i].X
		input_y := input_data[i].Y
		input_width := input_data[i].Width
		input_height := input_data[i].Height

		if input_width > 0 && input_height > 0{
			filtering_status:= check_intersection(main_data, input_data[i])
			if filtering_status == "valid"{			
				box := new(SavedData)
				box.X = input_x
				box.Y = input_y
				box.Width = input_width
				box.Height = input_height
				raw_time := time.Now()
				formatted := fmt.Sprintf(
					"%d-%02d-%02d %02d:%02d:%02d",
					raw_time.Year(), raw_time.Month(), raw_time.Day(),
					raw_time.Hour(), raw_time.Minute(), raw_time.Second())

				box.Time = formatted

				// file, _ := json.MarshalIndent(*box, " ", "")
				file, err_save := json.Marshal(*box)
				if err_save != nil{
					status_save = false
					log.Println("[INFO] Error saving rectangle coordinates")
				}else{
					f_slice = append(f_slice, string(file))
					status_save = true
				}

			}else{
				filtered_box := fmt.Sprintf(
					"%d,%d,%d,%d",
					input_x, input_y,
					input_width, input_height)
				log.Println("[INFO] This input recatangle ["+filtered_box+"] is filtered!")
			}
		}else{
			invalid_coordinate_box := fmt.Sprintf(
				"%d,%d,%d,%d",
				input_x, input_y,
				input_width, input_height)
			log.Println("[INFO] width or height of this recatangle ["+invalid_coordinate_box+"] is Invalid!")
		}

	}

	if status_save{
		for _, value := range f_slice {
			_, _ = datawriter.WriteString(value + ",\n")
			log.Println("[INFO] This input recatangle "+value+" is saved.")
		}
	}

	datawriter.Flush()
	Resp_file.Close()

}


// This function check if a key exists in a json data
func keyExists(key string, keys []string) bool {
    for _, k := range keys {
        if k == key {
            return true
        }
    }
    return false
}

// This function check if an id exists in a json data
func check_id(id int, invalid_idx []int) bool{
	for _, idx := range invalid_idx {
        if idx == id {
            return false
        }
    }
    return true

}


/*
	This function test object names of request data including
	{"main" and "input"}
	and also it checks the feild names
	including {x, y, width, height}
*/
func check_request_object_name(Ext_map map[string]interface{}) (
	string,
	MainData,
	[]InputData){

	var main_obj MainData
	var input_obj []InputData
	status_pass := "failed"

	for k, val := range Ext_map {
		if !keyExists(k, allowedExtKeys) {
			log.Println("[INFO] Disallowed key `"+k+"` in JSON")
			ret_str := "Invalid \"main\" or \"input\" tag name in request data"
			return ret_str, main_obj, input_obj
		}
		if k == "main"{
			for k_2, _ := range val.(map[string]interface{}){
				if !keyExists(k_2, allowedIntKeys) {
					log.Println("[INFO] Disallowed key `"+k_2+"` in \"main\" JSON request data!")
					ret_str2 := "Invalid \"main\" property name"
					return ret_str2, main_obj, input_obj
				}
			}
			num_main_feilds := reflect.ValueOf(val).Len()
			if num_main_feilds == Num_Feilds{
				mapstructure.Decode(val, &main_obj)
				status_pass = "passed"
			}else{
				ret_str3 := "Invalid number of feilds in \"main\" property is lower than 4"
				return ret_str3, main_obj, input_obj
			}
		}else if k == "input"{
			arr := val.([]interface{})
			len_input := len(arr)

			var invalid_idx []int
			for idx:= 0; idx < len_input; idx++{
				for k_2, _ := range val.([]interface{})[idx].(map[string]interface{}){
					if !keyExists(k_2, allowedIntKeys) {
						invalid_idx = append(invalid_idx, idx)
						log.Println("[Info] Input_index[",idx , "], Disallowed key `"+k_2+"` in \"input\" JSON request data!")
					}
				}
			}
			var filtered_arr []interface{}
			var num_input_feilds int
			for kk, elm := range arr {
				status_id := check_id(kk, invalid_idx)
				if status_id{
					num_input_feilds = reflect.ValueOf(elm).Len()
					if num_input_feilds == Num_Feilds{
						filtered_arr = append(filtered_arr, elm)
					}else{
						log_str4 := "Invalid number of feilds, one or more \"input\" property is lower than 4"
						log.Println("[INFO] "+log_str4)
					}
				}
			}

			mapstructure.Decode(filtered_arr, &input_obj)
		}
	}
	return status_pass, main_obj, input_obj
}


/* 
   This function handles POST requset
   Using this command to run it:
   curl -X POST -s localhost:8080 \
   -d '{"main": {"x": 0, "y": 0, "width": 10, "height": 10}, "input": [{"x": 4, "y": 5, "width": 2, "height": 2}]}'
*/
func handlePost(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        panic(err)
    }

	Ext_map := map[string]interface{}{}

	if err := json.Unmarshal(body, &Ext_map); err != nil {
		log.Println("[Info] Error unmatshalling json data = ", err)
	}

	valid_status_object, main_obj, input_obj := check_request_object_name(Ext_map)

	if valid_status_object == "passed"{
		parse_value(main_obj, input_obj)

	}else{
		str := " , Try entering a correct request format again:"
		log.Println("[INFO] " + valid_status_object+ str)
	}
}


func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", handlePost).Methods("POST")
	router.HandleFunc("/", handleGet).Methods("GET")

	err_server := http.ListenAndServe(":" + PortNum, router)

	if errors.Is(err_server, http.ErrServerClosed) {
		log.Println("[INFO] Server closed")
	} else if err_server != nil {
		log.Println("[INFO] Error starting server: ", err_server)
		os.Exit(1)
	}
}