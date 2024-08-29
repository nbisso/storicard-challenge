La idea es armar el procesamiento de los archivos por topico, de esta forma se podran
procesar en paralelo las transacciones de bulk y mantener estado y retry de las mismas

El archivo sera guardado para poder procesarlo luego (ya que si es muy pesado y queremos usar esto en un ui o llamar a otra api damos response rapido y comenzamos a procesar)