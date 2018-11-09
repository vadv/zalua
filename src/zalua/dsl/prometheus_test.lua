prometheus.listen(":1111")

-- simple gauge
gauge_one, err = prometheus_gauge.new({help = "this is help", name = "one_gauge"})
if err then error(err) end

gauge_one:add(10)
gauge_one:set(1)

-- vector gauge
gauge_vec_one, err = prometheus_gauge_labels.new({help = "this is help", name = "one_gauge_vec", labels = {"user"}})
if err then error(err) end

gauge_vec_two, err = prometheus_gauge_labels.new({help = "this is help", name = "one_gauge_vec", labels = {"user"}})
if err then error(err) end

gauge_vec_one:add({user = "user_1"}, 10)
gauge_vec_one:add({user = "user_2"}, 20)
gauge_vec_two:add({user = "user_2"}, 2)
gauge_vec_two:set({user = "user_3"}, 1)
gauge_vec_two:set({user = "user_3"}, 2)

-- simple counter
counter_one, err = prometheus_counter.new({help = "this is help", name = "counter_one"})
if err then error(err) end

counter_two, err = prometheus_counter.new({help = "this is help", name = "counter_one"})
if err then error(err) end

counter_one:inc()
counter_two:add(1000.2)


-- vector counter
counter_vec_one, err = prometheus_counter_lables.new({help = "this is help", name = "one_counter_vec", labels = {"user"}})
if err then error(err) end

counter_vec_two, err = prometheus_counter_lables.new({help = "this is help", name = "one_counter_vec", labels = {"user"}})
if err then error(err) end

counter_vec_one:add({user = "user_1"}, 10)
counter_vec_one:add({user = "user_2"}, 20)
counter_vec_two:add({user = "user_2"}, 2)

