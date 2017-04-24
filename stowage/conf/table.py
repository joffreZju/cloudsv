#! usr/bin/python
#coding=utf-8
#filetype=python
import sys,psycopg2

def alterTable():
    conn = psycopg2.connect(database="stowage",user="allsum",password="stowage@allsum,./",host="127.0.0.1",port="5432")
    print 'connect successful!'
    cursor=conn.cursor()
    usersql="alter table public.allsum_user alter column password drop default;"
    agentsql='''alter table agent 
                alter column name drop default,
                alter column uid drop default,
                alter column license_file drop default;'''
    accountsql = '''alter table account
                    alter column userid drop default,
                    alter column account_no drop default'''
    couponsql = '''alter table coupon
                    alter column number drop default,
                    alter column denomination drop default,
                    alter column verify_code drop default'''
    documentsql = '''alter table document
                    alter column doc_type drop default,
                    alter column uploader drop default,
                    alter column file_no drop default'''
    filesql = '''alter table file
                    alter column file_no drop default,
                    alter column name drop default,
                    alter column md5 drop default,
                    alter column uid drop default'''
    ordersql = '''alter table allsum_order
                    alter column order_no drop default,
                    alter column uid drop default,
                    alter column create_time drop default,
                    alter column price drop default'''
    billsql = '''alter table bill
                    alter column user_id drop default,
                    alter column time drop default,
                    alter column account_id drop default,
                    alter column money drop default'''
    fileTp = '''alter table file
                add column datax bytea,
                drop column data,
                rename column datax to data'''
    cursor.execute(usersql)
    cursor.execute(agentsql)
    cursor.execute(accountsql)
    cursor.execute(couponsql)
    cursor.execute(documentsql)
    cursor.execute(filesql)
    cursor.execute(ordersql)
    cursor.execute(billsql)
    #cursor.execute(fileTp)
    
    conn.commit()
    #cursor.close()
    conn.close()

if __name__ == "__main__":
    alterTable()
