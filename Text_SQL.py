from flask import Flask, request, jsonify
import os
from openai import OpenAI
from sqlalchemy import create_engine, text, MetaData, Table, Column, Integer, String
from sqlalchemy.exc import SQLAlchemyError
import csv
import json
from io import StringIO


# Initialize OpenAI API client
api_key = os.getenv("OPENAI_API_KEY")
client = OpenAI(api_key=api_key)

app = Flask(__name__)

# Function to convert text to SQL using the OpenAI API
def text_to_sql(text_query):
    prompt = f"Convert the following text to a SQL query:\n\n{text_query}\n\nSQL query:"
    
    try:
        response = client.chat.completions.create(
            model="gpt-3.5-turbo",  # or "gpt-4" if available
            messages=[
                {"role": "system", "content": "You are a helpful assistant that converts text to SQL queries."},
                {"role": "user", "content": prompt}
            ],
            max_tokens=2000,
            n=1,
            temperature=0.7,
        )

        # Extract the generated SQL query
        sql_query = response.choices[0].message.content
        return sql_query
    except Exception as e:
        return f"An error occurred: {str(e)}"

# Function to establish a dynamic connection to the database using SQLAlchemy
def connect_to_database(db_url):
    try:
        # Create a SQLAlchemy engine
        engine = create_engine(db_url)
        return engine
    except Exception as e:
        raise ConnectionError(f"Failed to connect to the database: {str(e)}")



# API route to handle text queries and return the SQL results

# API route to upload a CSV file and convert it to JSON
# API route to handle JSON file and insert its content into a database
@app.route('/api/upload_json', methods=['POST'])
def upload_json():
    # Expect JSON payload with 'data' and 'db_url'
    try:
        json_data = request.json.get('data')  # The JSON data provided in the request body
        db_url = request.json.get('db_url')   # The database URL

        if not json_data or not db_url:
            return jsonify({'error': 'Invalid input. Ensure both data and db_url are provided.'}), 400
        
        # Connect to the database
        engine = connect_to_database(db_url)

        # Create table dynamically based on JSON keys (assuming flat JSON structure)
        with engine.connect() as connection:
            metadata = MetaData()
            table_name = 'uploaded_json_data'

            # Extract column names from JSON keys (assuming all items have the same keys)
            sample_record = json_data[0]  # Take the first record to define the columns
            columns = [Column('id', Integer, primary_key=True, autoincrement=True)]
            for column_name in sample_record.keys():
                columns.append(Column(column_name, String))  # Assuming all columns are strings for simplicity
            
            # Create the table
            table = Table(table_name, metadata, *columns)
            metadata.create_all(engine)

            # Insert the JSON data into the table
            connection.execute(table.insert(), json_data)

            return jsonify({
                'message': 'JSON data inserted into the database successfully.',
                'table_name': table_name
            }), 200

    except SQLAlchemyError as e:
        return jsonify({'error': str(e)}), 500
    except Exception as e:
        return jsonify({'error': str(e)}), 400

    finally:
        engine.dispose()


if __name__ == '__main__':
    app.run(debug=True)
