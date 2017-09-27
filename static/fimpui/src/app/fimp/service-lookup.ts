const FIMP_SERVICE_LIST = [
    {"name":"basic","label":"Generic level","icon":""},
    {"name":"dev_sys","label":"Device system","icon":""},
    {"name":"out_bin_switch","label":"Output switch/relay","icon":""},
    {"name":"out_lvl_switch","label":"Output level switch/dimmer","icon":""},
    {"name":"meter_elec","label":"Electricity meter","icon":""},
    {"name":"meter_gas","label":"Gas meter","icon":""},
    {"name":"meter_water","label":"Water meter","icon":""},
    {"name":"sensor_temp","label":"Temperature sensor","icon":""},
    {"name":"sensor_lumin","label":"Luminance sensor","icon":""},
    {"name":"sensor_contact","label":"Open/close sensor","icon":""},
    {"name":"sensor_presence","label":"Presence detection sensor","icon":""},
    {"name":"alarm_fire","label":"Fire alarm","icon":""},
    {"name":"alarm_heat","label":"Heat alarm","icon":""},
    {"name":"alarm_burglar","label":"Intrusion alarm","icon":""},
    {"name":"battery","label":"Battery level","icon":""},
    {"name":"thermostat","label":"Thermostat","icon":""},
    {"name":"door_lock","label":"Doorlock","icon":""},
    {"name":"color_ctrl","label":"Color control","icon":""},
    {"name":"scene_ctrl","label":"Scene controller","icon":""},
    {"name":"fan_ctrl","label":"Fan speed and modes","icon":""}

];

export function getFimpServiceList() {
    return FIMP_SERVICE_LIST
}