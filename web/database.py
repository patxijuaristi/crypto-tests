import sqlite3

def get_all_sizes_data():
    conn = sqlite3.connect('../crypto_tests.db')
    cursor = conn.cursor()
    cursor.execute("SELECT * FROM key_sig_sizes_data")
    rows = cursor.fetchall()
    conn.close()

    data_list = []
    for row in rows:
        data_dict = {
            'id': row[0],
            'algorithm': row[1],
            'private_key': row[2],
            'public_key': row[3],
            'signature': row[4]
        }
        data_list.append(data_dict)
    
    return data_list

def delete_all_data():
    conn = sqlite3.connect('../crypto_tests.db')
    cursor = conn.cursor()
    cursor.execute("DELETE FROM key_sig_sizes_data")
    conn.commit()
    conn.close()