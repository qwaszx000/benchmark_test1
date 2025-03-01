import os
import os.path
import re
import csv
import pandas

def process_wrk_report(file):
    result = [0] * 6

    #skip 3 lines
    for _ in range(3):
        file.readline()

    line_data_list = []
    for data_cell in file.readline().split(" "):
        stripped = data_cell.strip()
        if stripped:
            line_data_list.append(stripped)
    
    avg_latency_str = line_data_list[1]
    max_latency_str = line_data_list[3]

    #skip 6 lines
    for _ in range(6):
        file.readline()

    total_requests_sent = 0
    for data_cell in file.readline().split(" "):
        stripped = data_cell.strip()
        if stripped:
            total_requests_sent = int(stripped)
            break

    line = file.readline().strip()

    avg_rps = 0.0
    errors_count = 0

    def _parse_avg_rps_line(_line):
        for data_cell in reversed(_line.split(" ")):
            stripped = data_cell.strip()
            if stripped:
                return float(stripped)

    #Errors line and/or avg rps line
    is_avg_parsed = False
    if line.startswith("Socket"):
        for data_cell in line.split(" "):
            if data_cell.isdigit():
                errors_count + int(data_cell)
    else:
        avg_rps = _parse_avg_rps_line(line)
        is_avg_parsed = True

    if not is_avg_parsed:
        avg_rps = _parse_avg_rps_line(file.readline().strip())

    result[0] = total_requests_sent
    result[1] = avg_rps

    result[2] = 0 #no data
    result[3] = pandas.Timedelta(avg_latency_str).total_seconds()
    result[4] = pandas.Timedelta(max_latency_str).total_seconds()

    result[5] = (total_requests_sent - errors_count) / total_requests_sent

    return result

def process_go_wrk_report(file):
    result = [0] * 6

    #skip 2 lines
    for _ in range(2):
        file.readline()

    line = file.readline()
    if line.startswith("Error"):
        return result
    
    total_requests_sent = 0
    for data_cell in line.split(" "):
        stripped = data_cell.strip()
        if stripped and stripped.isdigit():
            total_requests_sent = int(stripped)
            break

    #skip 2 lines
    for _ in range(2):
        file.readline()

    avg_rps = float(file.readline().strip().split("\t")[-1])

    #skip 1 line
    file.readline()

    min_latency_str = file.readline().strip().split("\t")[-1]
    avg_latency_str = file.readline().strip().split("\t")[-1]
    max_latency_str = file.readline().strip().split("\t")[-1]

    errors_count = int(file.readline().strip().split("\t")[-1])

    result[0] = total_requests_sent
    result[1] = avg_rps

    result[2] = pandas.Timedelta(min_latency_str).total_seconds()
    result[3] = pandas.Timedelta(avg_latency_str).total_seconds()
    result[4] = pandas.Timedelta(max_latency_str).total_seconds()

    result[5] = (total_requests_sent - errors_count) / total_requests_sent

    return result

def process_my_util_report(file):
    result = [0] * 6

    while True:
        line = file.readline()
        if line.endswith("Done\n"):
            break
    
    #total_requests_sent = int(file.readline().strip().split(" ")[-1])
    #skip 1 line
    file.readline()

    avg_rps = float(file.readline().strip().split(" ")[-1])
    success_rate = float(file.readline().strip().split(" ")[-1][:-1]) / 100

    #skip 1 line
    file.readline()

    min_latency_str = file.readline().strip().split(" ")[-1]
    max_latency_str = file.readline().strip().split(" ")[-1]
    avg_latency_str = file.readline().strip().split(" ")[-1]

    total_requests_sent = int(file.readline().strip().split(" ")[-1])

    result[0] = total_requests_sent
    result[1] = avg_rps

    result[2] = pandas.Timedelta(min_latency_str).total_seconds()
    result[3] = pandas.Timedelta(avg_latency_str).total_seconds()
    result[4] = pandas.Timedelta(max_latency_str).total_seconds()

    result[5] = success_rate

    return result

def process_file(file_path):
    filename = os.path.basename(file_path)
    technology, test_tool = filename[:-4].split("__")
    
    data = [technology, test_tool]
    with open(file_path) as report_file:
        if test_tool == "wrk":
            return data + process_wrk_report(report_file)
        elif test_tool == "go-wrk":
            return data + process_go_wrk_report(report_file)
        elif test_tool == "my_util":
            return data + process_my_util_report(report_file)
        else:
            raise ValueError("test_tool is not supported")

def main():
    with open("results.csv", "w") as csv_file:
        results_writer = csv.writer(csv_file)

        header_row = [
            "Technology",
            "Tool",

            "Total requests made",
            "AVG RPS",

            "Min latency",
            "AVG latency",
            "Max latency",

            "Success rate"
        ]
        results_writer.writerow(header_row)

        for entity in os.scandir("../results"):
            if entity.is_file():
                print(entity.name)

                csv_row_data = process_file(entity.path)
                results_writer.writerow(csv_row_data)

if __name__ == "__main__":
    main()