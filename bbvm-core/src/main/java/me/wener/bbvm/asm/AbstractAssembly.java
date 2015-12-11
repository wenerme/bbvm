package me.wener.bbvm.asm;

/**
 * @author wener
 * @since 15/12/11
 */
abstract class AbstractAssembly implements Assembly {
    Comment comment;

    public Comment getComment() {
        return comment;
    }

    public void setComment(Comment comment) {
        this.comment = comment;
    }

    protected String commentAssembly() {
        return comment != null ? " " + comment.toAssembly() : "";
    }

    public boolean hasComment() {
        return comment != null;
    }
}
