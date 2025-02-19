How do I know when reference types are automatically deferenced 


```
how does item.(tasks.Task) know to assign the value of item after type assertion and not both the value and a boolean (to know if it's of the type as this is a comma-ok pattern?)
data := &InitialData{
			Count: 1,
			Tasks: map[int]tasks.Task{
				1: item.(tasks.Task),
			},
		}
```