extends Control

# load the SIMPLE library
const SIMPLE = preload("res://bin/simple.gdns")
var data = SIMPLE.new()

func _on_Button_pressed():
	# data comes directly from the Go context
	$Label.text = "Data = " + data.get_data()
	data.HP += 1
	print("data.Blood is ", data.Blood)
	data.Blood += 1
