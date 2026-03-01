package planner

missing[action] {
  input.goal.trip.require_hotel
  not input.current.hotel_reserved
  action := "reserve_hotel"
}

missing[action] {
  input.goal.trip.require_dinner
  not input.current.dinner_reserved
  action := "reserve_dinner"
}
