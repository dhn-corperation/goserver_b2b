package databasepool

var resultBackProc string = `
CREATE PROCEDURE result_back_proc(p_interval int)
BEGIN
    DECLARE v_userid varchar(20);
    DECLARE v_ul_done INT DEFAULT FALSE;

    DECLARE user_list cursor for
            select distinct user_id from DHN_CLIENT_LIST dcl where dcl.use_flag  = 'Y';

    DECLARE CONTINUE HANDLER FOR NOT FOUND SET v_ul_done = TRUE;

	set @v_table_name = CONCAT( 'create table if not exists ', 'DHN_RESULT_' , DATE_FORMAT(now(), '%Y%m'), ' like DHN_RESULT')  ;
	
	prepare stmt from @v_table_name;
	execute stmt;
	deallocate prepare stmt;

    open user_list;
	ul_loop: 
	LOOP
        fetch user_list into v_userid;
        if v_ul_done then 
                close user_list;
                leave ul_loop;
        end if;

		set @insert_query = CONCAT( 'insert ignore into DHN_RESULT_BK_TEMP select * from DHN_RESULT where reg_dt <= DATE_FORMAT( date_sub( now(), INTERVAL  ',p_interval,' minute), ''%Y-%m-%d %T'') and userid =''', v_userid, ''' and sync = ''Y'' ');
		
		prepare stmt from @insert_query;
		execute stmt;
		deallocate prepare stmt;
	
		INS_BLOCK:
		begin
			DECLARE v_ins_done int Default false;
			DECLARE v_reg_dt varchar(20);
			DECLARE v_tbl_dt varchar(20);
		
			DECLARE c_ins_list_dt cursor for
  				select distinct concat(LAST_DAY(SUBSTR(reg_dt, 1,10)), ' 23:59:59') as reg_dt, concat (SUBSTR(reg_dt, 1,4),SUBSTR(reg_dt, 6,2)) as tbl_dt from DHN_RESULT_BK_TEMP where sync = 'Y' order by 1;

  			DECLARE CONTINUE HANDLER FOR NOT FOUND SET v_ins_done = TRUE;
  		
  			open c_ins_list_dt;
  			ins_loop:
  			loop
  				fetch c_ins_list_dt into v_reg_dt
  				                        ,v_tbl_dt;
  				if v_ins_done then
  					close c_ins_list_dt;
  				    leave ins_loop;
  				end if;
			
  			    set @ins_query = CONCAT('insert ignore into ', 'DHN_RESULT_', v_tbl_dt, ' select * from DHN_RESULT_BK_TEMP where userid = ''', v_userid , ''' and reg_dt <= ''', v_reg_dt, ''' ');
				prepare stmt from @ins_query;
				execute stmt;
				deallocate prepare stmt;	
			
				set @ud_query = CONCAT('update DHN_RESULT_BK_TEMP set sync=''C'' where userid = ''', v_userid , ''' and reg_dt <= ''', v_reg_dt, ''' ');
				prepare stmt from @ud_query;
				execute stmt;
				deallocate prepare stmt;	
			
  			end loop ins_loop;
  			
		end INS_BLOCK;
	
        PROC_BLOCK:
		begin
			DECLARE v_res_done INT DEFAULT FALSE;
			DECLARE v_userid varchar(20);
			DECLARE v_msgid  varchar(20);
			DECLARE v_reg_dt varchar(20);

			DECLARE c_result_list cursor for
				select userid, msgid , concat (SUBSTR(reg_dt, 1,4),SUBSTR(reg_dt, 6,2)) as reg_dt from DHN_RESULT_BK_TEMP where sync = 'C';
			
			DECLARE CONTINUE HANDLER FOR NOT FOUND SET v_res_done = TRUE;
			
			open c_result_list;
			rl_loop: 
			Loop
				fetch c_result_list into v_userid 
		  								,v_msgid 
										,v_reg_dt;
					if v_res_done then 
						close c_result_list;
						leave rl_loop;
					end if;
	
				delete from DHN_RESULT where userid = v_userid and msgid = v_msgid;
	
			END LOOP rl_loop;
		
			truncate table DHN_RESULT_BK_TEMP;
		end PROC_BLOCK;
	END LOOP ul_loop;
END
`

var resultSataProc string = `
CREATE PROCEDURE result_sata_proc(IN p_date VARCHAR(8))
BEGIN
    DECLARE v_userid VARCHAR(20);
    DECLARE v_ul_done INT DEFAULT FALSE;
    DECLARE v_table_name VARCHAR(50);
    DECLARE v_last_month VARCHAR(50);
    DECLARE v_pre_table_name VARCHAR(50);

    DECLARE user_list CURSOR FOR
        SELECT DISTINCT user_id FROM DHN_CLIENT_LIST dcl WHERE dcl.use_flag = 'Y';

    DECLARE CONTINUE HANDLER FOR NOT FOUND SET v_ul_done = TRUE;

    OPEN user_list;
    ul_loop: LOOP
        FETCH user_list INTO v_userid;
        IF v_ul_done THEN
            
            
            LEAVE ul_loop;
        END IF;

    	 SELECT v_userid AS userid;
        
        SET v_table_name = CONCAT('DHN_RESULT_', LEFT(p_date, 6));
        SET v_last_month = DATE_FORMAT(STR_TO_DATE(CONCAT(SUBSTRING(p_date, 1, 6), '01'), '%Y%m%d') - INTERVAL 1 MONTH, '%Y%m');
        SET v_pre_table_name = CONCAT('DHN_RESULT_', v_last_month);

        
        DROP TEMPORARY TABLE IF EXISTS tmp_result;

        
        SET @query = CONCAT('CREATE TEMPORARY TABLE tmp_result AS ',
                            'SELECT userid, ', 
                            'p_invoice AS depart, ',
                            'COUNT(1) AS send_cnt, ',
                            'SUBSTR(dr.reg_dt, 1, 10) AS send_date, ',
                            'SUM(CASE WHEN dr.s_code = ''0000'' AND dr.remark3 = 2 THEN 1 ELSE 0 END) AS ats_cnt, ',
                            'SUM(CASE WHEN dr.s_code <> ''0000'' AND dr.remark3 = 2 THEN 1 ELSE 0 END) AS ate_cnt, ',
                            'SUM(CASE WHEN dr.s_code = ''0000'' AND dr.remark3 = 1 AND dr.image_url IS NULL THEN 1 ELSE 0 END) AS fts_cnt, ',
                            'SUM(CASE WHEN dr.s_code <> ''0000'' AND dr.remark3 = 1 AND dr.image_url IS NULL THEN 1 ELSE 0 END) AS fte_cnt, ',
                            'SUM(CASE WHEN dr.s_code = ''0000'' AND dr.remark3 = 1 AND dr.image_url IS NOT NULL AND dr.wide = ''N'' THEN 1 ELSE 0 END) AS ftis_cnt, ',
                            'SUM(CASE WHEN dr.s_code <> ''0000'' AND dr.remark3 = 1 AND dr.image_url IS NOT NULL AND dr.wide = ''N'' THEN 1 ELSE 0 END) AS ftie_cnt, ',
                            'SUM(CASE WHEN dr.s_code = ''0000'' AND dr.remark3 = 1 AND dr.wide = ''Y'' THEN 1 ELSE 0 END) AS ftws_cnt, ',
                            'SUM(CASE WHEN dr.s_code <> ''0000'' AND dr.remark3 = 1 AND dr.wide = ''Y'' THEN 1 ELSE 0 END) AS ftwe_cnt, ',
                            'SUM(CASE WHEN dr.code = ''0000'' AND lower(dr.message_type) = ''ph'' AND lower(dr.sms_kind) = ''s'' THEN 1 ELSE 0 END) AS smss_cnt, ',
                            '0 AS smsd_cnt, ',
                            'SUM(CASE WHEN dr.code <> ''0000'' AND lower(dr.message_type) = ''ph'' AND lower(dr.sms_kind) = ''s'' THEN 1 ELSE 0 END) AS smse_cnt, ',
                            'SUM(CASE WHEN dr.code = ''0000'' AND lower(dr.message_type) = ''ph'' AND lower(dr.sms_kind) = ''l'' THEN 1 ELSE 0 END) AS lmss_cnt, ',
                            '0 AS lmsd_cnt, ',
                            'SUM(CASE WHEN dr.code <> ''0000'' AND lower(dr.message_type) = ''ph'' AND lower(dr.sms_kind) = ''l'' THEN 1 ELSE 0 END) AS lmse_cnt, ',
                            'SUM(CASE WHEN dr.code = ''0000'' AND lower(dr.message_type) = ''ph'' AND lower(dr.sms_kind) = ''m'' THEN 1 ELSE 0 END) AS mmss_cnt, ',
                            '0 AS mmsd_cnt, ',
                            'SUM(CASE WHEN dr.code <> ''0000'' AND lower(dr.message_type) = ''ph'' AND lower(dr.sms_kind) = ''m'' THEN 1 ELSE 0 END) AS mmse_cnt ',
                            'FROM ', v_pre_table_name, ' dr ',
                            'WHERE userid = ''', v_userid, ''' ',
                            'AND dr.reg_dt BETWEEN DATE_SUB(STR_TO_DATE(''', p_date, ''', ''%Y%m%d''), INTERVAL 5 DAY) AND DATE_ADD(STR_TO_DATE(''', p_date, ''', ''%Y%m%d''), INTERVAL 1 DAY) ',
                            'GROUP BY userid, p_invoice, SUBSTR(dr.reg_dt, 1, 10) ',
                            'UNION ALL ',
                            'SELECT userid, ', 
                            'p_invoice AS depart, ',
                            'COUNT(1) AS send_cnt, ',
                            'SUBSTR(dr.reg_dt, 1, 10) AS send_date, ',
                            'SUM(CASE WHEN dr.s_code = ''0000'' AND dr.remark3 = 2 THEN 1 ELSE 0 END) AS ats_cnt, ',
                            'SUM(CASE WHEN dr.s_code <> ''0000'' AND dr.remark3 = 2 THEN 1 ELSE 0 END) AS ate_cnt, ',
                            'SUM(CASE WHEN dr.s_code = ''0000'' AND dr.remark3 = 1 AND dr.image_url IS NULL THEN 1 ELSE 0 END) AS fts_cnt, ',
                            'SUM(CASE WHEN dr.s_code <> ''0000'' AND dr.remark3 = 1 AND dr.image_url IS NULL THEN 1 ELSE 0 END) AS fte_cnt, ',
                            'SUM(CASE WHEN dr.s_code = ''0000'' AND dr.remark3 = 1 AND dr.image_url IS NOT NULL AND dr.wide = ''N'' THEN 1 ELSE 0 END) AS ftis_cnt, ',
                            'SUM(CASE WHEN dr.s_code <> ''0000'' AND dr.remark3 = 1 AND dr.image_url IS NOT NULL AND dr.wide = ''N'' THEN 1 ELSE 0 END) AS ftie_cnt, ',
                            'SUM(CASE WHEN dr.s_code = ''0000'' AND dr.remark3 = 1 AND dr.wide = ''Y'' THEN 1 ELSE 0 END) AS ftws_cnt, ',
                            'SUM(CASE WHEN dr.s_code <> ''0000'' AND dr.remark3 = 1 AND dr.wide = ''Y'' THEN 1 ELSE 0 END) AS ftwe_cnt, ',
                            'SUM(CASE WHEN dr.code = ''0000'' AND lower(dr.message_type) = ''ph'' AND lower(dr.sms_kind) = ''s'' THEN 1 ELSE 0 END) AS smss_cnt, ',
                            '(SELECT COUNT(1) FROM DHN_RESULT a WHERE userid=dr.userid AND a.p_invoice=dr.p_invoice AND a.result = ''P'' AND lower(a.sms_kind) = ''s'' AND SUBSTR(a.res_dt, 1, 10) = SUBSTR(dr.res_dt, 1, 10)) AS smsd_cnt, ',
                            'SUM(CASE WHEN dr.code <> ''0000'' AND lower(dr.message_type) = ''ph'' AND lower(dr.sms_kind) = ''s'' THEN 1 ELSE 0 END) AS smse_cnt, ',
                            'SUM(CASE WHEN dr.code = ''0000'' AND lower(dr.message_type) = ''ph'' AND lower(dr.sms_kind) = ''l'' THEN 1 ELSE 0 END) AS lmss_cnt, ',
                            '(SELECT COUNT(1) FROM DHN_RESULT a WHERE userid=dr.userid AND a.p_invoice=dr.p_invoice AND a.result = ''P'' AND lower(a.sms_kind) = ''l'' AND SUBSTR(a.res_dt, 1, 10) = SUBSTR(dr.res_dt, 1, 10)) AS lmsd_cnt, ',
                            'SUM(CASE WHEN dr.code <> ''0000'' AND lower(dr.message_type) = ''ph'' AND lower(dr.sms_kind) = ''l'' THEN 1 ELSE 0 END) AS lmse_cnt, ',
                            'SUM(CASE WHEN dr.code = ''0000'' AND lower(dr.message_type) = ''ph'' AND lower(dr.sms_kind) = ''m'' THEN 1 ELSE 0 END) AS mmss_cnt, ',
                            '(SELECT COUNT(1) FROM DHN_RESULT a WHERE userid=dr.userid AND a.p_invoice=dr.p_invoice AND a.result = ''P'' AND lower(a.sms_kind) = ''m'' AND SUBSTR(a.res_dt, 1, 10) = SUBSTR(dr.res_dt, 1, 10)) AS mmsd_cnt, ',
                            'SUM(CASE WHEN dr.code <> ''0000'' AND lower(dr.message_type) = ''ph'' AND lower(dr.sms_kind) = ''m'' THEN 1 ELSE 0 END) AS mmse_cnt ',
                            'FROM ', v_table_name, ' dr ',
                            'WHERE userid = ''', v_userid, ''' AND dr.sync = ''Y'' ',
                            'AND dr.reg_dt BETWEEN DATE_SUB(STR_TO_DATE(''', p_date, ''', ''%Y%m%d''), INTERVAL 5 DAY) AND DATE_ADD(STR_TO_DATE(''', p_date, ''', ''%Y%m%d''), INTERVAL 1 DAY) ',
                            'GROUP BY userid, p_invoice, SUBSTR(dr.reg_dt, 1, 10)');

        
        PREPARE stmt FROM @query;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;

        
        BEGIN
            DECLARE v_res_done INT DEFAULT FALSE;
            DECLARE v_ruserid VARCHAR(20);
            DECLARE v_depart VARCHAR(24);
            DECLARE v_send_date VARCHAR(20);
            DECLARE v_send_cnt INT;
            DECLARE v_ats_cnt INT;
            DECLARE v_ate_cnt INT;
            DECLARE v_fts_cnt INT;
            DECLARE v_fte_cnt INT;
            DECLARE v_ftis_cnt INT;
            DECLARE v_ftie_cnt INT;
            DECLARE v_ftws_cnt INT;
            DECLARE v_ftwe_cnt INT;
            DECLARE v_smss_cnt INT;
            DECLARE v_smsd_cnt INT;
            DECLARE v_smse_cnt INT;
            DECLARE v_lmss_cnt INT;
            DECLARE v_lmsd_cnt INT;
            DECLARE v_lmse_cnt INT;
            DECLARE v_mmss_cnt INT;
            DECLARE v_mmsd_cnt INT;
            DECLARE v_mmse_cnt INT;

            DECLARE res CURSOR FOR
                SELECT userid, depart, send_date, sum(send_cnt) AS send_cnt, sum(ats_cnt) AS ats_cnt, sum(ate_cnt) AS ate_cnt, 
                	     sum(fts_cnt) AS fts_cnt, sum(fte_cnt) AS fte_cnt, sum(ftis_cnt) AS ftis_cnt, sum(ftie_cnt) AS ftie_cnt, 
                	     sum(ftws_cnt) AS ftws_cnt, sum(ftwe_cnt) AS ftwe_cnt, sum(smss_cnt) AS smss_cnt, sum(smsd_cnt) AS smsd_cnt, 
                	     sum(smse_cnt) AS smse_cnt, sum(lmss_cnt) AS lmss_cnt, sum(lmsd_cnt) AS lmsd_cnt, 
                	     sum(lmse_cnt) AS lmse_cnt, sum(mmss_cnt) AS mmss_cnt, sum(mmsd_cnt) AS mmsd_cnt, sum(mmse_cnt) AS mmse_cnt
                FROM tmp_result GROUP BY userid, depart, send_date;

            DECLARE CONTINUE HANDLER FOR NOT FOUND SET v_res_done = TRUE;

            OPEN res;
            res_loop: LOOP
                FETCH res INTO v_ruserid, v_depart, v_send_date, v_send_cnt, v_ats_cnt, v_ate_cnt, 
                              v_fts_cnt, v_fte_cnt, v_ftis_cnt, v_ftie_cnt, v_ftws_cnt, v_ftwe_cnt, 
                              v_smss_cnt, v_smsd_cnt, v_smse_cnt, v_lmss_cnt, v_lmsd_cnt, v_lmse_cnt, v_mmss_cnt, v_mmsd_cnt, v_mmse_cnt;
                IF v_res_done THEN
                	  
                    LEAVE res_loop;
                END IF;

                
                
                DELETE FROM DHN_RESULT_STA WHERE send_date = v_send_date AND userid = v_ruserid AND depart = v_depart;

                
                INSERT INTO DHN_RESULT_STA
                (send_date, userid, depart, send_cnt, ats_cnt, ate_cnt, fts_cnt, fte_cnt, ftis_cnt, ftie_cnt, ftws_cnt, ftwe_cnt, smss_cnt, smsd_cnt, smse_cnt, lmss_cnt, lmsd_cnt, lmse_cnt, mmss_cnt, mmsd_cnt, mmse_cnt)
                VALUES(v_send_date, v_ruserid, v_depart, v_send_cnt, v_ats_cnt, v_ate_cnt, v_fts_cnt, v_fte_cnt, v_ftis_cnt, v_ftie_cnt, v_ftws_cnt, v_ftwe_cnt, v_smss_cnt, v_smsd_cnt, v_smse_cnt, v_lmss_cnt, v_lmsd_cnt, v_lmse_cnt, v_mmss_cnt, v_mmsd_cnt, v_mmse_cnt);

                COMMIT;
            END LOOP res_loop;
            CLOSE res;

            
            DROP TEMPORARY TABLE IF EXISTS tmp_result;
        END;
    END LOOP ul_loop;
    CLOSE user_list;
    COMMIT;
END
`

var removeWs string = `
CREATE FUNCTION remove_ws(P_MSG LONGTEXT CHARACTER set utf8mb4) RETURNS longtext CHARSET utf8mb4
BEGIN
	DECLARE v_msg LONGTEXT  CHARACTER set utf8;
	DECLARE v_ohc varchar(20);
	DECLARE v_ds  varchar(10) CHARACTER set utf8;
	DECLARE v_done INT DEFAULT FALSE;

	DECLARE sc cursor for
	  select origin_hex_code
	        ,dest_str
	    from SPECIAL_CHARACTER sc 
	   where enabled  = 'Y';

  	DECLARE CONTINUE HANDLER FOR NOT FOUND SET v_done = TRUE;

	set v_msg = p_msg;  


	open sc;
	sc_loop: LOOP
		fetch sc
		  into v_ohc
		      ,v_ds;
		     
		set v_msg = replace(v_msg, unhex(v_ohc),ifnull(v_ds,''));
		if v_done then 
			leave sc_loop;
		end if;
	END LOOP;
	
	close sc;

    return v_msg;
END
`