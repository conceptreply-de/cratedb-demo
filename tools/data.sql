

insert into vehicles (id, vin, description, use) values ('1', '5YJSA1DG9DFP14705', 'Honda Accord 2018', 'Route 15 Bus');

insert into vehicles (id, vin, description, use) values ('2', '5YJSA1DG9DFP14890', 'Toyota Camry 2020', 'Route 17 Bus');


insert into vehicle_measurements (
  vehicle_id,
  measurement_type,
  measurement_time,
  measurement_value
) VALUES ( '2', 'STATE_OF_CHARGE', '2024-01-08 04:05:06'::timestamp, 5 ),
         ( '2', 'STATE_OF_CHARGE', '2024-01-08 04:06:06'::timestamp, 14 ),
         ( '2', 'STATE_OF_CHARGE', '2024-01-08 04:20:06'::timestamp, 25 );
