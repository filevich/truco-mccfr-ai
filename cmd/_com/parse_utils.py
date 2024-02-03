import datetime

parse_date = lambda s: datetime.datetime.strptime(s, '%Y/%m/%d %H:%M:%S')

def progress(
    current,
    goal,
    start:datetime.datetime,
    last_commit:datetime.datetime,
):
    prog = current / goal
    delta = last_commit - start
    eta_total = delta.total_seconds() / prog
    eta_total = datetime.timedelta(seconds=round(eta_total))
    prog = round(prog, 3)
    return prog, delta, eta_total-delta, eta_total