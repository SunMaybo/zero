package template

const JavaDaoTemplate = `
package cn.xunhou.grpc.{{.ServiceName}}.dao;

import cn.xunhou.cloud.dao.xhjdbc.XbbRepository;
import cn.xunhou.grpc.xhportal.entity.{{.TableUpper}}Entity;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.jdbc.core.namedparam.NamedParameterJdbcTemplate;
import org.springframework.stereotype.Repository;

/**
 * @author fico
 */
@Repository
public class {{.TableUpper}}Repository extends XbbRepository<{{.TableUpper}}Entity> {
    private final Logger LOGGER = LoggerFactory.getLogger(this.getClass());

    public MenusRepository(NamedParameterJdbcTemplate jdbcTemplate) {
        super(jdbcTemplate);

    }
}

`
