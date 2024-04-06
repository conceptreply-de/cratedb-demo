
insert into vehicle_measurements (
  vehicle_id,
  measurement_type,
  measurement_time,
  measurement_value
) VALUES ( '2', 'STATE_OF_CHARGE', '2024-01-08 04:05:06'::timestamp, 5 ),
         ( '2', 'STATE_OF_CHARGE', '2024-01-08 04:06:06'::timestamp, 14 ),
         ( '2', 'STATE_OF_CHARGE', '2024-01-08 04:20:06'::timestamp, 25 );


INSERT INTO vehicles (id, vin, description, use)
VALUES
    ('1', '5YJSA1DG9DFP14705', 'Volvo 700 Electric', 'Urban transportation'),
    ('2', '5YJSA1DG9DFP14890', 'Mercedes-Benz Sprinter', 'Tourism'),
    ('3', 'JM1BL1SG2A1123456', 'MAN Lion''s City', 'City commuting'),
    ('4', 'WDDUG8FB7EA000001', 'Scania Citywide', 'Intercity transport'),
    ('5', '2HKRM3H33EH501234', 'Neoplan Tourliner', 'Long-distance travel'),
    ('6', 'WAUZZZ8E3AA500001', 'Iveco Crossway', 'Public transportation'),
    ('7', '1N4AA5AP7DC800001', 'VDL Citea', 'Airport shuttle'),
    ('8', 'VF1JL14B3D3030147', 'Solaris Urbino', 'School bus'),
    ('9', '5XXGN4A70CG800001', 'Volvo 9700', 'Coach travel'),
    ('10', 'JTDDR32T7Y0021221', 'Van Hool TX', 'Cross-country touring'),
    ('11', 'KNDMB5C14H6325948', 'BYD K9', 'Electric bus fleet'),
    ('12', '1N6AD0ERXCC445429', 'Setra S 431 DT', 'Luxury touring'),
    ('13', '2HNYD2H23DH500001', 'Volkswagen Crafter', 'Delivery service'),
    ('14', 'JH4DB1640CS003207', 'Mitsubishi Fuso', 'Cargo transport'),
    ('15', '1N4BA41E94C800001', 'Temsa Prestij', 'Corporate shuttle'),
    ('16', 'WAUZZZ8U2DA100001', 'Mercedes-Benz Travego', 'VIP transport'),
    ('17', '1G1FP22T0NL108799', 'Solaris InterUrbino', 'Regional travel'),
    ('18', '1GNEC13Z72R240760', 'Ford Transit', 'Mobile office'),
    ('19', 'WBABW53424PL46319', 'Optare Solo', 'Local community service'),
    ('20', 'JM1BK323551200001', 'Isuzu Erga', 'Special needs transportation'),
    ('21', 'WAUZZZ8E3AA500002', 'Iveco Crossway', 'Public transportation'),
    ('22', '1N4AA5AP7DC800002', 'VDL Citea', 'Airport shuttle'),
    ('23', 'VF1JL14B3D3030148', 'Solaris Urbino', 'School bus'),
    ('24', '5XXGN4A70CG800002', 'Volvo 9700', 'Coach travel'),
    ('25', 'JTDDR32T7Y0021222', 'Van Hool TX', 'Cross-country touring'),
    ('26', 'KNDMB5C14H6325949', 'BYD K9', 'Electric bus fleet'),
    ('27', '1N6AD0ERXCC445428', 'Setra S 431 DT', 'Luxury touring'),
    ('28', '2HNYD2H23DH500002', 'Volkswagen Crafter', 'Delivery service'),
    ('29', 'JH4DB1640CS003208', 'Mitsubishi Fuso', 'Cargo transport'),
    ('30', '1N4BA41E94C800002', 'Temsa Prestij', 'Corporate shuttle');

