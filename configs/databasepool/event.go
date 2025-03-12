package databasepool

var evtResultBackProc string = `
CREATE EVENT evt_result_back_proc
ON SCHEDULE EVERY 1 DAY
STARTS '2023-10-28 00:20:00.000'
ON COMPLETION NOT PRESERVE
ENABLE
DO CALL result_back_proc(20)
`

var evtResultSataProc string = `
CREATE EVENT evt_result_sata_proc
ON SCHEDULE EVERY 1 DAY
STARTS '2024-04-28 00:50:00.000'
ON COMPLETION NOT PRESERVE
ENABLE
DO BEGIN
    DECLARE previous_date VARCHAR(8);
    SET previous_date = DATE_FORMAT(DATE_SUB(CURDATE(), INTERVAL 1 DAY), '%Y%m%d');
    CALL result_sata_proc(previous_date);
END
`

var evtRemoveReception string = `
CREATE EVENT evt_remove_reception
ON SCHEDULE EVERY 1 DAY
STARTS '2024-04-28 00:20:00.000'
ON COMPLETION NOT PRESERVE
ENABLE
DO BEGIN
    delete from DHN_RECEPTION where insert_date < date_sub(now(), interval 15 DAY);
END
`