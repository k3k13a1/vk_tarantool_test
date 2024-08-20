box.cfg{
    listen = 3301,
    log_level = 5
}

-- Создаем пространство, если его не существует
local space = box.schema.space.create('vk_test', {
    if_not_exists = true
})

-- Добавляем первичный индекс
space:create_index('primary', {
    type = 'tree',
    parts = {1, 'string'}, -- Первый индекс по строковому ключу (например, ключи из JSON)
    if_not_exists = true
})

-- Дополнительный индекс (например, по значению)
space:create_index('value_index', {
    type = 'tree',
    parts = {2, 'scalar'}, -- Индекс по второму полю, который может быть любого типа
    if_not_exists = true
})

-- Создаем пользователя и задаем ему права
box.schema.user.create('testuser', {password = 'pass', if_not_exists = true})
box.schema.user.grant('testuser', 'read,write,execute', 'universe', nil, {if_not_exists = true})

print("Tarantool is ready to use!")
