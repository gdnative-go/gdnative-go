extends Control

# load the SIMPLE library
const SIMPLE = preload("res://bin/simple.gdns")
onready var data = SIMPLE.new()

func _on_Button_pressed():
	$Label.text = "Data = " + data.get_data()
