subinclude("//build_defs/wings")

wings_lib(
  name = "day",
  config = "wings.json",
  srcs = ["day.wings"],
  outs = {
    "src/cockpit/wings": "Day.ts",
    "src/turbine/wings": "day.go",
  },
  deps = [
    "//:wings_config",
    ":place",
  ],
  visibility = ["PUBLIC"],
)

wings_lib(
  name = "new_user",
  config = "wings.json",
  srcs = ["new_user.wings"],
  outs = {
    "src/cockpit/wings": "NewUser.ts",
    "src/turbine/wings": "newuser.go",
  },
  deps = [
    "//:wings_config",
  ],
  visibility = ["PUBLIC"],
)

wings_lib(
  name = "place",
  config = "wings.json",
  srcs = ["place.wings"],
  outs = {
    "src/cockpit/wings": "Place.ts",
    "src/turbine/wings": "place.go",
  },
  deps = [
    "//:wings_config",
  ],
  visibility = ["PUBLIC"],
)

wings_lib(
  name = "trip",
  config = "wings.json",
  srcs = ["trip.wings"],
  outs = {
    "src/cockpit/wings": "Trip.ts",
    "src/turbine/wings": "trip.go",
  },
  deps = [
    "//:wings_config",
    ":day",
    "//src/wings/enum:city"
  ],
  visibility = ["PUBLIC"],
)

wings_lib(
  name = "user",
  config = "wings.json",
  srcs = ["user.wings"],
  outs = {
    "src/cockpit/wings": "User.ts",
    "src/turbine/wings": "user.go",
  },
  deps = [
    "//:wings_config",
  ],
  visibility = ["PUBLIC"],
)
